Name:           bx
Version:        1.5.3
Release:        1%{?dist}
Summary:        Command-Line Tool for 1C-Bitrix Module Development

License:        MIT
URL:            https://github.com/pixel365/%{name}
Source0:        %{name}_v%{version}_linux_amd64.tar.gz

%global debug_package %{nil}

Requires:       glibc

%description
BX is a command-line tool for developers working on 1C-Bitrix platform modules.
It allows you to declaratively define all stages of project build,
as well as validate the module configuration and deploy the final distribution.
Build configurations are versioned alongside the project,
ensuring consistency and traceability of changes throughout the development process.

%prep
%setup -q -c -T
tar -xzf %{SOURCE0}

%build
# nothing to build

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_licensedir}/%{name}
mkdir -p %{buildroot}%{_docdir}/%{name}

install -m0755 bx %{buildroot}%{_bindir}/%{name}
install -m0644 LICENSE %{buildroot}%{_licensedir}/%{name}/LICENSE
install -m0644 README.md %{buildroot}%{_docdir}/%{name}/README.md

%files
%{_bindir}/%{name}
%license %{_licensedir}/%{name}/LICENSE
%doc %{_docdir}/%{name}/README.md

%changelog
* Tue May 20 2025 Ruslan Semagin <pixel.365.24@gmail.com> - 1.5.3-1
- Initial binary release from GitHub
