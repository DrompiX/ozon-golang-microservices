set -e

sed -i "$ a include_if_exists = '/etc/postgresql/postgres.conf'" /var/lib/postgresql/data/postgresql.conf
