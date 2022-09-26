$cmd = $args[0]
$params = $args[1]

$db_host="127.0.0.1"
$db_port="3003"
$db_user="root"
$db_pass="root"
$db_name="stellar_db"

$db_addr = "mysql://{2}:{3}@tcp({0}:{1})/{4}?parseTime=true" -f ($db_host, $db_port, $db_user, $db_pass, $db_name)

if (($cmd -eq "up") -or ($cmd -eq "down" )) {
    if ($params -ne "") {
        migrate -source file://./db/migrations -database $db_addr $cmd $params
    } else {
        migrate -source file://./db/migrations -database $db_addr $cmd
    }
} elseif ($cmd -eq "drop") {
    migrate -source file://./db/migrations -database $db_addr $cmd
} elseif ($cmd -eq "new") {
    migrate create -ext sql -dir ./db/migrations $params
}