-- name: UpsertNode :exec
INSERT INTO nodes (node_id, ip, port, cpus, cpu_arch)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (node_id) DO UPDATE
    SET ip = excluded.ip,
        port = excluded.port,
        cpus = excluded.cpus,
        cpu_arch = excluded.cpu_arch;

-- name: GetNodes :many
SELECT * FROM nodes;

-- name: DeleteNode :exec
DELETE FROM nodes WHERE node_id = ?;
