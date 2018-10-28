#!/bin/sh

echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
echo "root:test"|chpasswd
/entrypoint.sh "$@"