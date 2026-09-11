# Test Migration

Jalankan migration product pada database `be_evindo`:

```powershell
psql -U postgres -h localhost -p 5432 -d be_evindo -f .\migrations\test\001_product_test.sql
```


