# Test Migration

Jalankan migration pada database `be_evindo`:

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\001_product_test.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\002_create_users.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\003_create_suppliers.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\004_create_warehouses.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\005_create_inventories.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\006_create_purchase_requests.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\007_create_purchase_orders.sql
```


