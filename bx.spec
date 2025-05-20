Name:           bx
Version:        1.5.3
Release:        1%{?dist}
Summary:        Command-Line Tool for 1C-Bitrix Module Development

License:        MIT
URL:            https://github.com/pixel365/%{name}
Source0:        https://github.com/pixel365/%{name}/releases/download/v%{version}/%{name}_%{version}_linux_amd64.tar.gz

Requires:       glibc

%description
BX is a command-line tool for developers working on 1C-Bitrix platform modules.
It allows you to declaratively define all stages of project build,
as well as validate the module configuration and deploy the final distribution.
Build configurations are versioned alongside the project,
ensuring consistency and traceability of changes throughout the development process.

%prep
mkdir -p build
cp %{SOURCE0} build/
cd build
tar -xzf %{name}_%{version}_linux_amd64.tar.gz

%build
# nothing to build

%install
install -D -m0755 build/%{name} %{buildroot}%{_bindir}/%{name}
install -D -m0644 build/LICENSE %{buildroot}%{_licensedir}/%{name}/LICENSE
install -D -m0644 build/README.md %{buildroot}%{_docdir}/%{name}/README.md

%files
%{_bindir}/%{name}
%license %{_licensedir}/%{name}/LICENSE
%doc %{_docdir}/%{name}/README.md

%changelog
* Mon May 20 2025 Ruslan Semagin <you@example.com> - 1.5.3-1
- Initial binary release from GitHub
