<?php

date_default_timezone_set('Asia/Shanghai');

if (!function_exists('getallheaders')) {
    function getallheaders()
    {
        $headers = '';
        foreach ($_SERVER as $name => $value) {
            if (substr($name, 0, 5) == 'HTTP_') {
                $headers[str_replace(' ', '-', ucwords(strtolower(str_replace('_', ' ', substr($name, 5)))))] = $value;
            }
        }
        return $headers;
    }
}

$headers = getallheaders();

$log = "\nNew request: " . date('Y-m-d H:i:s') . ":\n" .
    var_export(
        array(
            'method' => isset($_SERVER['REQUEST_METHOD']) ? $_SERVER['REQUEST_METHOD'] : '',
            'headers' => $headers,
            'cookies' => $_COOKIE,
            'get' => $_GET,
            'raw_post' => file_get_contents('php://input'),
            'post' => $_POST
        ),
        true
    );

file_put_contents(
    'ab.log',
    $log,
    FILE_APPEND
);

if (!(isset($headers['Test']) && $headers['Test'] == 'test') ||
    !(isset($headers['Test1']) && $headers['Test1'] == 'test1')
) {
    http_response_code(403);
    exit;
}

if (!(isset($_COOKIE['test']) && $_COOKIE['test'] == 'test') ||
    !(isset($_COOKIE['test1']) && $_COOKIE['test1'] == 'test1')
) {
    http_response_code(403);
    exit;
}

if (!(isset($_POST['test']) && $_POST['test'] == 'test') ||
    !(isset($_POST['test1']) && $_POST['test1'] == 'test1')
) {
    http_response_code(403);
    exit;
}