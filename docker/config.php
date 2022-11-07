<?php  // Moodle configuration file

unset($CFG);
global $CFG;
$CFG = new stdClass();

$CFG->dbtype    = 'mysqli';
$CFG->dblibrary = 'native';
$CFG->dbhost    = getenv('DB_HOST');
$CFG->dbname    = getenv('DB_NAME');
$CFG->dbuser    = getenv('DB_USER');
$CFG->dbpass    = getenv('DB_PASS');
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

$CFG->directorypermissions = 0777;

require_once(__DIR__ . '/lib/setup.php');
