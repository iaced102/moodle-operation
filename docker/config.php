<?php  // Moodle configuration file

unset($CFG);
global $CFG;
$CFG = new stdClass();

$CFG->dbtype    = 'mysqli';
$CFG->dblibrary = 'native';
$CFG->dbhost    = '45.124.94.128';
$CFG->dbname    = 'moodle';
// $CFG->dbuser    = 'duy';
// $CFG->dbpass    = 'TY3Q5qd26baAkQdcV7gHtR8';
$CFG->dbuser    = 'duy';
$CFG->dbpass    = '4Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL';
$CFG->prefix    = 'mdl_';
$CFG->dboptions = array (
  'dbpersist' => 0,
  'dbport' => '',
  'dbsocket' => '',
  // 'dbcollation' => 'utf8mb4_unicode_ci',
  'dbcollation' => 'utf8_general_ci',
);

$CFG->wwwroot   = 'http://'.$_SERVER['SERVER_NAME'];
$CFG->dataroot  = '/var/www/moodledata';
$CFG->dirroot   = '/var/www/moodle';
$CFG->admin     = 'admin';

$CFG->directorypermissions = 0777;

require_once(__DIR__ . '/lib/setup.php');

// There is no php closing tag in this file,
// it is intentional because it prevents trailing whitespace problems!
