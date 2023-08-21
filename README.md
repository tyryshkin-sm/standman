# About

CLI-утилита для поднятия стенда из QEMU-виртуальных машин. Основана на библиотеке libvirt.
На данный момент сборка образов и настройка в них сети производятся в ручном режиме.

## Build

```bash
make build
```

## Usage

В директории с конфигурационным файлом standman.yaml.

```bash
standman up
standman down
```

## Appendix

```bash
sudo apt install libvirt-dev

https://libvirt.org/apps.html
https://libvirt.org/formatdomain.html

virsh list --inactive --name | xargs -r -n 1 virsh undefine
```
