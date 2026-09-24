# How to setup the arcade in Kiosk mode

This guide assume we’re using Linux Bazzite with KDE, Wayland, and Plasma

## Launcher Scripts

Use three scripts: one outer script to create the DBus session, one compositor script to start KWin and XWayland, and one app-session script that KWin launches after the display environment is ready.

### `/var/home/arcade/session-launcher.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

exec dbus-run-session -- /var/home/arcade/session-launcher-inner.sh
```

### `/var/home/arcade/session-launcher-inner.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

export XDG_SESSION_TYPE=wayland
export XDG_CURRENT_DESKTOP=KDE
export XDG_SESSION_DESKTOP=KDE
export DESKTOP_SESSION=arcade
export QT_QPA_PLATFORM=wayland
export GDK_BACKEND=wayland

exec kwin_wayland \
  --drm \
  --xwayland \
  --no-lockscreen \
  --exit-with-session=/var/home/arcade/app-session.sh
```

Steam still requires XWayland. The `--xwayland` option starts KWin's rootless XWayland server and adds its `DISPLAY` and `XAUTHORITY` values to the environment of the session command.

Do not start the back-end before `kwin_wayland`. Processes started before KWin cannot inherit the display environment created by the compositor, so Steam exits with `Unable to open X11 display`.

### `/var/home/arcade/app-session.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

BACKEND=/var/home/arcade/arcade-champion/back-end/back-end
FRONTEND=/var/home/arcade/arcade-champion/tauri/target/release/arcade-champion

"$BACKEND" &
backend_pid=$!

cleanup() {
  kill "$backend_pid" 2>/dev/null || true
  wait "$backend_pid" 2>/dev/null || true
}

trap cleanup EXIT INT TERM

"$FRONTEND"
```

KWin starts this script only after the Wayland and XWayland displays are available. Both the back-end and front-end therefore inherit the same `WAYLAND_DISPLAY`, `DISPLAY`, `XAUTHORITY`, and DBus session values. Games launched by the back-end inherit them in turn.

Make all scripts executable:

```bash
chmod 0755 /var/home/arcade/session-launcher.sh
chmod 0755 /var/home/arcade/session-launcher-inner.sh
chmod 0755 /var/home/arcade/app-session.sh
```

This is intentionally a **minimal KWin kiosk session**, not a complete Plasma desktop.

## Desktop entry

Needs to be bundled into a RPM package because it should be added in /usr, which is immutable on Bazzite.

Create a Fedora development container:

```bash
toolbox create --container rpm-builder
toolbox enter --container rpm-builder
```

Inside the container, install the build tool:

```bash
sudo dnf install -y rpm-build
```

Still inside the toolbox container, run:

```bash
mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
```

Create the session file:

```bash
cat > ~/rpmbuild/SOURCES/arcade.desktop <<'EOF'

[Desktop Entry]
Type=Application
Name=Arcade Champion
Comment=Arcade Champion kiosk session
Exec=/var/home/arcade/session-launcher.sh
TryExec=/var/home/arcade/session-launcher.sh
DesktopNames=KDE
EOF
```

Create the package specification:

```bash
cat > ~/rpmbuild/SPECS/arcade-champion-session.spec <<'EOF'
Name:           arcade-champion-session
Version:        1.0
Release:        1%{?dist}
Summary:        Arcade Champion Wayland session
License:        MIT
BuildArch:      noarch
Source0:        arcade.desktop

%description
Plasma Login Wayland session entry for Arcade Champion.

%install
install -Dpm 0644 %{SOURCE0} \
%{buildroot}%{_datadir}/wayland-sessions/arcade.desktop

%files
%{_datadir}/wayland-sessions/arcade.desktop
EOF
```

Build it:

```bash
rpmbuild -bb ~/rpmbuild/SPECS/arcade-champion-session.spec
```

Exit the toolbox:

```bash
exit
```

Your home directory is shared with Toolbox, so the RPM created in `~/rpmbuild/` is visible from the Bazzite host as well.

Layer it into the Bazzite deployment:

```bash
sudo rpm-ostree install \
  ~/rpmbuild/RPMS/noarch/arcade-champion-session*
```

Then reboot into the new deployment:

```bash
systemctl reboot
```

After the reboot:

```bash
ls -l /usr/share/wayland-sessions/
```

Expected:

```text
arcade.desktop
plasma.desktop
```

Cleanup:

```bash
toolbox rm -f rpm-builder
```

## Plasma Login configuration

Reboot, try to log in using the Arcade Champion session. If everything works well, go to System Settings \-> Connexion screen and enable autologin in this session (do not enable automated relogin).
