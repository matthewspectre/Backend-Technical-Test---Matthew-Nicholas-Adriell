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

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\008_create_purchase_order_items.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\009_update_purchase_order_status_constraint.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\010_create_goods_receipts.sql
```

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\011_create_goods_receipt_items.sql
```


