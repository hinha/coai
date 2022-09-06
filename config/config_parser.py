import yaml

if __name__ == '__main__':
    import os
    file_path = os.path.join(os.getcwd(), "config", "config.example.yml")
    with open(file_path) as fileobj:
        cfg = yaml.safe_load(fileobj)

    assert cfg['server']['secret_key'] == "EXAMPLE"
    cfg['server']['secret_key'] = os.getenv("APP_SECRET_KEY")

    for jwt in cfg['jwt']:
        if jwt.endswith("key"):
            assert cfg['jwt'][jwt] == "EXAMPLE"
            cfg['jwt'][jwt] = os.getenv(f"JWT_{jwt.upper()}")

    for db in cfg['database']['drivers']['mysql']:
        if db == "driver" or db == "port":
            continue
        assert cfg['database']['drivers']['mysql'][db] == "EXAMPLE"
        if db.startswith("db"):
            cfg['database']['drivers']['mysql'][db] = os.getenv(f"{db.upper()}")
        else:
            cfg['database']['drivers']['mysql'][db] = os.getenv(f"DB_{db.upper()}")

    file_path = os.path.join(os.getcwd(), "config", "config.yml")
    with open(file_path, 'w') as f:
        yaml.dump(cfg, f)

