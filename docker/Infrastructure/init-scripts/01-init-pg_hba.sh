#!/bin/bash
set -e

cat > "$PGDATA/pg_hba.conf" << EOF
# TYPE  DATABASE        USER            ADDRESS                 METHOD

# Local connections
local   all             all                                     trust

# IPv4 local connections:
host    all             all             127.0.0.1/32            trust
host    all             all             0.0.0.0/0               trust

# IPv6 local connections:
host    all             all             ::1/128                 trust
host    all             all             ::0/0                   trust
EOF

# Reload PostgreSQL configuration
pg_ctl reload -D "$PGDATA" || true
