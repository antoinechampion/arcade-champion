Summary:        Arcade Champion plymouth boot theme
Name:           custom-plymouth-theme
Version:        1.0.0
Release:        1%{?dist}
License:        GPLv2+
Group:          System Environment/Base
URL:            https://github.com/kkklemennn/rpm-custom-plymouth-theme
Source0:        custom-plymouth-theme-1.0.0.tar.gz

BuildArch:      noarch
Requires:       plymouth

%define themedir %{_datadir}/plymouth/themes/arcade

%description
The %{name} package contains the Arcade Champion boot theme for plymouth.

%prep
%autosetup -p1

%install
install -d -m 755 %{buildroot}%{themedir}
install -m 644 -p arcade.plymouth %{buildroot}%{themedir}/arcade.plymouth
install -m 644 -p *.png %{buildroot}%{themedir}/

%files
%dir %{themedir}
%{themedir}/*

%changelog
* Wed Mar 06 2024 Klemen Klemar <klemen.klemar@hotmail.com> - 1.0.0
- Created the theme
