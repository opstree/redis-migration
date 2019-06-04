# Redis Migration
The main idea of writing this migration utility was to provide a easy and convenient migration.

## Dependencies

The list of dependencies are not quite long but yes we do have some dependencies.

#### System Dependencies
- [X] **python**

#### Python Dependencies
- [X] **redis**

Don't worry we have taken care the python dependencies in [requirments.txt](./requirments.txt)

## Usage
```shell
pip install -r requirments.txt
```

Update redis connection url in [migration.py](./migration.py)

```python
OLD = "old-redis.opstree.com"
NEW = "new-redis.opstree.com"
```

That's it now run the python utility

```shell
python migration.py
```
