import stun

def test_nat_type():
    servers = [
        ("stun.l.google.com", 19302),
        ("stun1.l.google.com", 19302)
    ]

    results = []
    for host, port in servers:
        nat_type, external_ip, external_port = stun.get_ip_info(stun_host=host, stun_port=port)
        print(f"STUN {host}:{port} => {external_ip}:{external_port} ({nat_type})")
        results.append((external_ip, external_port))

    symmetric = len(set([r[1] for r in results])) > 1
    if symmetric:
        print("\n🚫 Detected SYMMETRIC NAT (bad for punching).")
    else:
        print("\n✅ Non-symmetric NAT (good for punching).")

test_nat_type()
