# How to install Python packages via pip
> pip is the package installer for Python. You can us pip to install packages
from the `Python Package Index` and other indexes. [read more](https://pypi.org/project/pip/)

## Requirements for installing Packages
* Ensure `Python` is installed and you can run `Python` from command line.
  ```shell
  # the latest version of Python is 3.x
  python3 --version

  # Python 3.6.9
  ```

* Ensure `pip3` is installed and you can run `pip3` from command line.
  ```shell
  pip3 --version

  # pip 9.0.1 from /usr/lib/python3/dist-packages (python 3.6)
  ```

* If `pip3` is not installed, then install it
  ```shell
  sudo apt-get update
  sudo apt-get install -y python3-pip
  ```
* Ensure pip3, setuptools, and wheel are up to Update
  ```shell
  python -m pip install --upgrade pip3 setuptools wheel
  ```

## Usage
> we use `pip` in below examples, you can change to `pip3` if you are using the latest version 3.x

### Installing package
- To install the latest version of **"SomeProject"**:
  ```shell
  pip install "SomeProject"
  ```
- To install a specific Version:
  ```shell
  pip install "SomeProject==1.4"
  ```
- To install greater than or equal to one version and less than another:
  ```shell
  pip install "SomeProject>=1,<2"
  ```
- To install a version that’s “compatible” with a certain version:
  ```shell
  pip install "SomeProject~=1.4.2"
  ```
### Upgrading packages
* Upgrade an already installed **"SomeProject"** to the latest:
  ```shell
  pip install --upgrade SomeProject
  ```
### Installing to the User site
* To install packages that are isolated to the current user, use the --user flag:
  ```shell
  pip install --user SomeProject
  ```
### Installing packages from requirement files
* List of packages in `requirements.txt`
  ```
  pkg1
  pkg2==2.1.8
  pkg3>=1.3.0
  pkg4>=1.0,<=2.0
  ```
* To install a list of requirements specific in a `requirements.txt`
  ```shell
  pip install -r requirements.txt
  ```
### Installing from VCS
* To install a project from VCS in 'editable' mode:
  ```shell
  pip install -e git+https://git.repo/some_pkg.git#egg=SomeProject          # from git
  pip install -e hg+https://hg.repo/some_pkg#egg=SomeProject                # from mercurial
  pip install -e svn+svn://svn.repo/some_pkg/trunk/#egg=SomeProject         # from svn
  pip install -e git+https://git.repo/some_pkg.git@feature#egg=SomeProject  # from a branch
  ```
### Installing from other indexes
* To install from an alternative indexes
  ```shell
  pip install --index-url http://my.package.repo/simple/ SomeProject
  ```
* To Search an additional index during install:
  ```shell
  pip install --extra-index-url http://my.package.repo/simple SomeProject
  ```
### Installing from a local src tree
* Installing from local src in `Development` mode, i.e. in such a way that the project appears to be installed, but yet is still editable from the src tree.
  ```shell
  pip install -e <path>
  ```  
* install normally from source
  ```shell
  pip install <path>
  ```
### Installing from local archives
* to install a particular source archive file:
  ```shell
  pip install ./downloads/SomeProject-1.0.4.tar.gz
  ```
* to install from a local directory containing archives:
  ```shell
  pip install --no-index --find-links=file:///local/dir/ SomeProject
  pip install --no-index --find-links=/local/dir/ SomeProject
  pip install --no-index --find-links=relative/dir/ SomeProject
  ```
