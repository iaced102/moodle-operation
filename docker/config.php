<?php  // Moodle configuration file

unset($CFG);
global $CFG;
$CFG = new stdClass();

$CFG->dbtype    = 'mariadb';
$CFG->dblibrary = 'native';
$CFG->dbhost    = getenv('DB_HOST');
$CFG->dbname    = getenv('DB_NAME');
$CFG->dbuser    = getenv('DB_USER');
$CFG->dbpass    = getenv('DB_PASS');
$CFG->theme		= getenv('THEME');
$CFG->prefix    = 'mdl_';
$CFG->dboptions = array (
  'dbpersist' => 0,
  'dbport' => '',
  'dbsocket' => '',
  'dbcollation' => 'utf8_general_ci',
);

$CFG->wwwroot   = 'http://'.$_SERVER['SERVER_NAME'];
$CFG->dataroot  = '/var/www/moodledata';
$CFG->dirroot   = '/var/www/moodle';
$CFG->admin     = 'admin';

$CFG->session_handler_class = '\core\session\redis';
$CFG->session_redis_host = getenv('REDIS_HOST');
$CFG->session_redis_port = getenv('REDIS_PORT');
$CFG->session_redis_auth = getenv('REDIS_PASS');
$CFG->session_redis_acquire_lock_timeout = getenv('REDIS_LOCK_TIMEOUT');
$CFG->session_redis_lock_expire = getenv('REDIS_LOCK_EXPIRE');
$CFG->session_redis_lock_retry = getenv('REDIS_LOCK_RETRY');


$CFG->directorypermissions = 0777;

require_once(__DIR__ . '/lib/setup.php');
