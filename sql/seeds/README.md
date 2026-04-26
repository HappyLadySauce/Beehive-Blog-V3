# Seeds

Seed SQL is development-only data. Run it explicitly after schema migrations when
the local database has been recreated.

种子 SQL 只用于开发数据。删库重刷后，先执行 schema migration，再显式执行 seed。

```powershell
.\sql\migrate.ps1
.\sql\seed.ps1
```

Current seeds:

- `v3/identity/001_dev_admin.sql`: creates or repairs the local admin account
  `admin@beehive.local / Admin@123456`.

Add seed SQL here for:

- owner bootstrap account
- default platform config
- development-only sample content
