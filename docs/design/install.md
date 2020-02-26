# 인스톨러 디자인 문서
## 매개변수화
매개변수화는 특정 변수를 치환해서 yaml 파일을 생성할 수 있게 해줍니다. 예를 들어 아래와 같은 yaml 파일이 있습니다.

```yaml
kind: Pod
metadata:
  name: foo
```

사용자는 Pod 이름을 foo가 아니라 bar로 만들고 싶을 수 있습니다. 이 때 우리는 이름을 아래와 같이 매개변수화 할 수 있습니다.

```yaml
kind: Pod
metadata:
  name: {{ pod_name }}
```

이제 Pod을 생성할 때는 sed나 curl과 같은 명령어를 통해 `{{ pod_name }}`을 찾아서 원하는 값으로 치환할 수 있습니다.

```shell
sed '{{ pod_name }}|bar|g' | kubectl apply -f -
```

우리가 관리하는 yaml 파일은 아주 많지만 수정되는 부분은 아주 많지 않습니다. 예를 들면 이미지 이름(레파지토리)등은 변경될 수 있는 값입니다.
하지만 kind등은 자주 변경되지 않을 것입니다. 따라서 이렇게 자주 변경되는 부분에 대해서만 매개변수화를 진행하자는 아이디어가 있었습니다.

### Helm Chart
우리는 매개변수화가 가지는 장단점을 파악하기 위해 실제로 Helm Chart를 적용해보았습니다. Helm Chart는 잘 작동했습니다. 하지만 다음과 같은 몇 가지 문제점을 찾을 수 있었습니다.
- Helm에 의해 생성된 yaml은 파일로 존재하지 않습니다.
  - Helm은 내부적으로 치환을 거치고 해당 yaml을 k8s에 적용합니다. 이 과정에서 yaml 파일이 특정 위치에 저장되지 않습니다. 이는 해당 클러스터를 yaml로 선언적으로 관리하기 어렵다는 것을 뜻합니다. 클러스터에 대한 yaml 파일은 클러스터의 모든 정보를 담고 있기 때문에 yaml 파일을 직접 보고 관리할 수 있는 것은 관리자에게 중요합니다.
- Helm 패키지 자체의 정보를 저장할 수 없습니다.
  - Infra As a Code(IaC)는 해당 인프라를 완전한 코드로 저장하고 있습니다. 따라서 해당 인프라가 소실되더라도 언제나 해당 코드를 통해 인프라를 다시 복구할 수 있습니다. 또한 인프라를 코드로 관리하기 때문에 git의 코드 관련 기능도 사용할 수 있습니다. 하지만 Helm을 사용하는 경우 IAC에 어려움이 따릅니다. 예를 들어 Helm으로 여러 패키지를 배포한 환경의 경우 그 환경의 정보(즉 어떤 패키지가 배포되었는가)를 코드로 관리하기가 어렵습니다.

## [parameterization-pitfalls](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/architecture/declarative-application-management.md#parameterization-pitfalls)
Kubernetes에서는 기본적으로 매개변수화에 대해서 부정적인 입장을 취하고 있습니다. 그 이유는 아래와 같습니다.

- 어떤 변수를 매개변수화 해야하나?
  - 자주 사용되는 변수를 매개변수화 한다고 했지만, 자주 사용되는 변수의 기준이 애매합니다. 처음에는 최소한의 변수만 매개변수화를 진행하더라도 이후 기능이 늘어감에 따라 더 많은 변수를 포함하게 됩니다. 그리고 언젠가는 모든 변수를 매개변수화하게 되어 매개변수화의 의미가 사라집니다.
- 디버깅의 어려움
  - 템플릿된 파일은 디버깅이 어렵습니다. 가장 큰 이유는 매개변수화를 진행한 이후의 yaml 파일이 특정 위치에 존재하는 것이 아니기에 차이점을 알아보기 어렵다는 문제입니다.

위와 같은 시행착오를 겪은 뒤 가장 중점적으로 둔 가치는 다음 세 가지입니다.
- 매개변수화는 없다: 매개변수화를 위한 코딩은 낭비적이고 매개변수화 자체 또한 그렇습니다.
- Kuberentes의 기능을 활용합니다: kubectl apply는 업그레이드를 포함해서 설계되었습니다. 따라서 이런 기능을 활용할 수 있어야 합니다.
- Yaml 파일이 명시적으로 존재합니다: IaC와 디버깅의 용이함을 위해서입니다.

### [Kustomize](https://kustomize.io/)
Kustomize는 Kubernetes Native한 방법으로 환경을 관리합니다. Kustomize에서는 크게 base와 overlay 영역이 존재합니다. base영역은 수정되지 않은 yaml 파일을 가지고 있습니다. 반면 overlay 영역에서는 base 영역에서 수정하고 싶은 파트만을 가지고 있습니다. 

```yaml
# base.yaml
kind: Pod
metadata:
  name: foo

# overlay.yaml
kind: Pod
metadata:
  name: bar
```

이제 이 yaml을 Kustomize를 통해 적용하면 base가 가진 모든 특성을 포함하되, overlay에 명시되어 있는 부분에 대해서는 overlay의 것을 따릅니다. 이 방법을 사용하면 다음과 같은 장점이 있습니다.

- 매개변수화를 위한 별다른 작업이 필요하지 않습니다. 치환하고 싶은 값이 있으면 단순히 overlay.yaml에 해당 파트를 추가해주면 됩니다.
- Yaml 파일이 명시적으로 존재합니다.
- 위 방법을 사용한 뒤 kubectl apply 커맨드를 사용할 수 있습니다.
- 디렉토리별로 다양한 환경을 관리할 수 있습니다. 예를 들어 3replica 환경, minikube 환경 등을 따로 만들어서 관리할 수 있습니다.

## 인스톨 순서
TODO
