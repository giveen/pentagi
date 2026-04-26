# How To Reset PentAGI Admin Password

This runbook is for local deployments that use Docker Compose in this repository.

## Scope

- Account: admin@pentagi.com
- Database container: pgvector
- Database name: pentagidb

## Option A: Emergency Access Reset (sets password to admin)

Use this only to regain access quickly, then change the password immediately.

From repository root:

    docker exec -i pgvector psql -U postgres -d pentagidb -c "UPDATE users SET password='\$2a\$10\$deVOk0o1nYRHpaVXjIcyCuRmaHvtoMN/2RUT7w5XbZTeiWKEbXx9q', password_change_required=false WHERE mail='admin@pentagi.com';"

Verify the row:

    docker exec -i pgvector psql -U postgres -d pentagidb -c "SELECT mail, status, password_change_required FROM users WHERE mail='admin@pentagi.com';"

Then login with:

- Email: admin@pentagi.com
- Password: admin

## Option B: Set a New Strong Password Directly in DB

1. Generate a bcrypt hash for your new password on your host:

    python3 - <<'PY'
    import getpass
    import bcrypt
    pw = getpass.getpass('New admin password: ').encode()
    print(bcrypt.hashpw(pw, bcrypt.gensalt(rounds=10)).decode())
    PY

If bcrypt is missing, install it first:

    python3 -m pip install bcrypt

2. Copy the generated hash and update the admin user:

    docker exec -i pgvector psql -U postgres -d pentagidb -c "UPDATE users SET password='<PASTE_BCRYPT_HASH_HERE>', password_change_required=false WHERE mail='admin@pentagi.com';"

3. Verify:

    docker exec -i pgvector psql -U postgres -d pentagidb -c "SELECT mail, status, password_change_required FROM users WHERE mail='admin@pentagi.com';"

## Password Policy

Use a password that meets these requirements:

- Minimum 12 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character
- Do not use common weak passwords

## Troubleshooting

If the SQL update affects 0 rows, confirm the account exists:

    docker exec -i pgvector psql -U postgres -d pentagidb -c "SELECT id, mail, status FROM users WHERE mail='admin@pentagi.com';"

If psql cannot connect, verify containers:

    docker ps --format "table {{.Names}}\t{{.Status}}"

Ensure pgvector is running and healthy before retrying.
