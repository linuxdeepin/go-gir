#add 2020-04-26
%global _unpackaged_files_terminate_build 0
%global with_debug 1

%if 0%{?with_debug}
%global debug_package   %{nil}
%endif

Name:           go-gir-generator
Version:        2.0.5
Release:        1
Summary:        Generate static golang bindings for GObject
License:        GPLv3
URL:            ssh://gerrit.uniontech.com:29418/go-gir-generator
Source0:        %{name}-%{version}.orig.tar.xz
#Patch0:         SettingsBackendLike.patch
# https://cr.deepin.io/#/c/go-gir-generator/+/41653/
#Patch1:         %{name}_build-with-glib2.patch

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required

# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires:  compiler(go-compiler)
BuildRequires:  pkgconfig(gobject-introspection-1.0)
BuildRequires:  pkgconfig(gudev-1.0)
Provides:       golang(gir/gobject-2.0)
Provides:       golang(gir/gio-2.0)
Provides:       golang(gir/glib-2.0)
Provides:       golang(gir/gudev-1.0)
Provides:       pkgconfig(github.com/linuxdeepin/go-gir/gio-2.0)
Provides:       deepin-gir-generator

%description
Generate static golang bindings for GObject

%prep
%setup -q -n %{name}-%{version}

#GIO_VER=$(v=$(rpm -q --qf %{RPMTAG_VERSION} gobject-introspection); echo ${v//./})
#if [ $GIO_VER -ge 1521 ]; then
# Our gobject-introspection is too new
# https://cr.deepin.io/#/c/16880/
#%patch0 -p1
#%patch1 -p1
#fi

%build
export GOPATH="%{gopath}"
%make_build

%install
%make_install

%files
%doc README.md
%license LICENSE
%{_bindir}/gir-generator
%{gopath}/src/github.com/linuxdeepin/go-gir/

%changelog
* Thu Mar 23 2021 uoser <uoser@uniontech.com> - 2.0.5-1
- Update to 2.0.5
