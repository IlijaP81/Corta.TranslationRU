// Code generated by statik. DO NOT EDIT.

// Package statik contains static assets.
package mysql

import (
	"github.com/rakyll/statik/fs"
)

func Data() string {
	return "PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1a\x00	\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE `crm_content` (\n `id` bigint(20) unsigned NOT NULL,\n `module_id` bigint(20) unsigned NOT NULL,\n `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n `updated_at` datetime DEFAULT NULL,\n `deleted_at` datetime DEFAULT NULL,\n PRIMARY KEY (`id`,`module_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `crm_content_column` (\n `content_id` bigint(20) NOT NULL,\n `column_name` varchar(255) NOT NULL,\n `column_value` text NOT NULL,\n PRIMARY KEY (`content_id`,`column_name`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `crm_field` (\n `field_type` varchar(16) NOT NULL COMMENT 'Short field type (string, boolean,...)',\n `field_name` varchar(255) NOT NULL COMMENT 'Description of field contents',\n `field_template` varchar(255) NOT NULL COMMENT 'HTML template file for field',\n PRIMARY KEY (`field_type`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `crm_module` (\n `id` bigint(20) unsigned NOT NULL,\n `name` varchar(64) NOT NULL COMMENT 'The name of the module',\n `json` json NOT NULL COMMENT 'List of field definitions for the module',\n `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n `updated_at` datetime DEFAULT NULL,\n `deleted_at` datetime DEFAULT NULL,\n PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `crm_module_form` (\n `module_id` bigint(20) unsigned NOT NULL,\n `place` tinyint(3) unsigned NOT NULL,\n `kind` varchar(64) NOT NULL COMMENT 'The type of the form input field',\n `name` varchar(64) NOT NULL COMMENT 'The name of the field in the form',\n `label` varchar(255) NOT NULL COMMENT 'The label of the form input',\n `help_text` text NOT NULL COMMENT 'Help text',\n `default_value` text NOT NULL COMMENT 'Default value',\n `max_length` int(10) unsigned NOT NULL COMMENT 'Maximum input length',\n `is_private` tinyint(1) NOT NULL COMMENT 'Contains personal/sensitive data?',\n PRIMARY KEY (`module_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE `crm_page` (\n `id` bigint(20) unsigned NOT NULL COMMENT 'Page ID',\n `self_id` bigint(20) unsigned NOT NULL COMMENT 'Parent Page ID',\n `module_id` bigint(20) unsigned NOT NULL COMMENT 'Module ID (optional)',\n `title` varchar(255) NOT NULL COMMENT 'Title (required)',\n `description` text NOT NULL COMMENT 'Description',\n `blocks` json NOT NULL COMMENT 'JSON array of blocks for the page',\n `visible` tinyint(4) NOT NULL COMMENT 'Is page visible in navigation?',\n `weight` int(11) NOT NULL COMMENT 'Order for navigation',\n PRIMARY KEY (`id`) USING BTREE,\n KEY `module_id` (`module_id`),\n KEY `self_id` (`self_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nPK\x07\x08\xac\xe8\x19\x1d\x12\n\x00\x00\x12\n\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00%\x00	\x0020180704080001.crm_fields-data.up.sqlUT\x05\x00\x01\x80Cm8INSERT INTO `crm_field` VALUES ('bool','Boolean value (yes / no)','');\nINSERT INTO `crm_field` VALUES ('email','E-mail input','');\nINSERT INTO `crm_field` VALUES ('enum','Single option picker','');\nINSERT INTO `crm_field` VALUES ('hidden','Hidden value','');\nINSERT INTO `crm_field` VALUES ('stamp','Date/time input','');\nINSERT INTO `crm_field` VALUES ('text','Text input','');\nINSERT INTO `crm_field` VALUES ('textarea','Text input (multi-line)','');\nPK\x07\x08f\x18\x1e\x84\xc5\x01\x00\x00\xc5\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00+\x00	\x0020181109133134.crm_content-ownership.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_content` ADD `user_id` BIGINT UNSIGNED NOT NULL AFTER `module_id`, ADD INDEX (`user_id`);\nPK\x07\x08\xeb!\x81\xc2k\x00\x00\x00k\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00.\x00	\x0020181109193047.crm_fields-related_types.up.sqlUT\x05\x00\x01\x80Cm8INSERT INTO `crm_field` (`field_type`, `field_name`, `field_template`) VALUES ('related', 'Related content', ''), ('related_multi', 'Related content (multiple)', '');PK\x07\x08:.\xfb8\xa6\x00\x00\x00\xa6\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x000\x00	\x0020181125122152.add_multiple_relationships.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE `crm_content_links` (\n `content_id` bigint(20) unsigned NOT NULL,\n `column_name` varchar(255) NOT NULL,\n `rel_content_id` bigint(20) unsigned NOT NULL,\n PRIMARY KEY (`content_id`,`column_name`,`rel_content_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;PK\x07\x08\xee\x12\x15	\x05\x01\x00\x00\x05\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00D\x00	\x0020181125132142.add_required_and_visible_to_module_form_fields.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_module_form` ADD `is_required` TINYINT(1) NOT NULL AFTER `is_private`, ADD `is_visible` TINYINT(1) NOT NULL AFTER `is_required`;PK\x07\x08\xa5q c\x91\x00\x00\x00\x91\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x005\x00	\x0020181202163130.fix-crm-module-form-primary-key.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_module_form` DROP PRIMARY KEY, ADD PRIMARY KEY(`module_id`, `place`);\nPK\x07\x08\xd9\xd4i\xe3W\x00\x00\x00W\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x000\x00	\x0020181204123650.add-crm-content-json-field.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_content` ADD `json` json DEFAULT NULL COMMENT 'Content in JSON format.' AFTER `user_id`;\nPK\x07\x08\"\x96\xd6pj\x00\x00\x00j\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x004\x00	\x0020181204155326.add-crm-module-form-json-field.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_module_form` ADD `json` JSON NOT NULL COMMENT 'Options in JSON format.' AFTER `kind`;PK\x07\x08\xb7\x93\xd4\xf6f\x00\x00\x00f\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00+\x00	\x0020181216214630.crm-content-to-record.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_content` RENAME TO `crm_record`;\nALTER TABLE `crm_record` MODIFY COLUMN `json` json DEFAULT NULL COMMENT 'Records in JSON format.';\n\nALTER TABLE `crm_content_column` RENAME TO `crm_record_column`;\nALTER TABLE `crm_record_column` CHANGE COLUMN `content_id` `record_id` bigint(20);\n\nALTER TABLE `crm_content_links` RENAME TO `crm_record_links`;\nALTER TABLE `crm_record_links` CHANGE COLUMN `content_id` `record_id` bigint(20) unsigned;\nALTER TABLE `crm_record_links` CHANGE COLUMN `rel_content_id` `rel_record_id` bigint(20) unsigned;\nPK\x07\x08mA\xa8\x1e&\x02\x00\x00&\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00$\x00	\x0020181217100000.add-charts-tbl.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE `crm_chart` (\n `id`         BIGINT(20)  UNSIGNED NOT NULL,\n `name`       VARCHAR(64)          NOT NULL COMMENT 'The name of the chart',\n `config`     JSON                 NOT NULL COMMENT 'Chart & reporting configuration',\n\n `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,\n `updated_at` DATETIME                      DEFAULT NULL,\n `deleted_at` DATETIME                      DEFAULT NULL,\n\n PRIMARY KEY (`id`)\n\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\nPK\x07\x08\xcf\xc6g\xf6\xe4\x01\x00\x00\xe4\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00#\x00	\x0020181224122301.rem-crm_field.up.sqlUT\x05\x00\x01\x80Cm8DROP TABLE `crm_field`;\nPK\x07\x08\xae \xfd2\x18\x00\x00\x00\x18\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00&\x00	\x0020190108100000.add-triggers-tbl.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE `crm_trigger` (\n `id`         BIGINT(20)  UNSIGNED NOT NULL,\n `name`       VARCHAR(64)          NOT NULL COMMENT 'The name of the trigger',\n `enabled`    BOOLEAN              NOT NULL COMMENT 'Trigger enabled?',\n `actions`    TEXT                 NOT NULL COMMENT 'All actions that trigger it',\n `source`     TEXT                 NOT NULL COMMENT 'Trigger source',\n `rel_module` BIGINT(20)  UNSIGNED     NULL COMMENT 'Primary module',\n\n `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,\n `updated_at` DATETIME                      DEFAULT NULL,\n `deleted_at` DATETIME                      DEFAULT NULL,\n\n PRIMARY KEY (`id`)\n\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\nPK\x07\x08+\xad\xb7\xed\xb8\x02\x00\x00\xb8\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00/\x00	\x0020190110175924.rem-crm-record-json-field.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_record` DROP COLUMN `json`;\nPK\x07\x08\x94#\xb9\x99-\x00\x00\x00-\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x008\x00	\x0020190114072000.cleanup-record-tables-and-multival.up.sqlUT\x05\x00\x01\x80Cm8-- No more links, we'll handle this through ref field on crm_record_value tbl\nDROP TABLE IF EXISTS `crm_record_links`;\n\n-- Not columns, values\nALTER TABLE `crm_record_column` RENAME TO `crm_record_value`;\n\n-- Simplify names\nALTER TABLE `crm_record_value` CHANGE COLUMN `column_name`  `name`  VARCHAR(64);\nALTER TABLE `crm_record_value` CHANGE COLUMN `column_value` `value` TEXT;\n\n-- Add reference\nALTER TABLE `crm_record_value` ADD  COLUMN `ref` BIGINT UNSIGNED DEFAULT 0 NOT NULL;\nALTER TABLE `crm_record_value` ADD  COLUMN `deleted_at` datetime DEFAULT NULL;\nALTER TABLE `crm_record_value` ADD  COLUMN `place` INT UNSIGNED DEFAULT 0 NOT NULL;\nALTER TABLE `crm_record_value` DROP PRIMARY KEY, ADD PRIMARY KEY(`record_id`, `name`, `place`);\nCREATE INDEX crm_record_value_ref ON crm_record_value (ref);\n\n\n-- We want this as a real field\nALTER TABLE `crm_module_form`  ADD  COLUMN `is_multi` TINYINT(1) NOT NULL;\n\n-- This will be handled through meta(json) fieldd\nALTER TABLE `crm_module_form`  DROP COLUMN `help_text`;\nALTER TABLE `crm_module_form`  DROP COLUMN `max_length`;\nALTER TABLE `crm_module_form`  DROP COLUMN `default_Value`;\nPK\x07\x08\x04]{\x1fo\x04\x00\x00o\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'\x00	\x0020190121132408.record-updated-by.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE `crm_record` CHANGE COLUMN `user_id`  `owned_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE `crm_record` ADD COLUMN `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE `crm_record` ADD COLUMN `updated_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE `crm_record` ADD COLUMN `deleted_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;\nPK\x07\x08\x8f\xecQt`\x01\x00\x00`\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00migrations.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE IF NOT EXISTS `migrations` (\n `project` varchar(16) NOT NULL COMMENT 'sam, crm, ...',\n `filename` varchar(255) NOT NULL COMMENT 'yyyymmddHHMMSS.sql',\n `statement_index` int(11) NOT NULL COMMENT 'Statement number from SQL file',\n `status` text NOT NULL COMMENT 'ok or full error message',\n PRIMARY KEY (`project`,`filename`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nPK\x07\x089S\x05%x\x01\x00\x00x\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x06\x00	\x00new.shUT\x05\x00\x01\x80Cm8#!/bin/bash\ntouch $(date +%Y%m%d%H%M%S).up.sql\nPK\x07\x08\xc1h\xf1\xfb/\x00\x00\x00/\x00\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xac\xe8\x19\x1d\x12\n\x00\x00\x12\n\x00\x00\x1a\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(f\x18\x1e\x84\xc5\x01\x00\x00\xc5\x01\x00\x00%\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81c\n\x00\x0020180704080001.crm_fields-data.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xeb!\x81\xc2k\x00\x00\x00k\x00\x00\x00+\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x84\x0c\x00\x0020181109133134.crm_content-ownership.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(:.\xfb8\xa6\x00\x00\x00\xa6\x00\x00\x00.\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81Q\x0d\x00\x0020181109193047.crm_fields-related_types.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xee\x12\x15	\x05\x01\x00\x00\x05\x01\x00\x000\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\\\x0e\x00\x0020181125122152.add_multiple_relationships.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xa5q c\x91\x00\x00\x00\x91\x00\x00\x00D\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xc8\x0f\x00\x0020181125132142.add_required_and_visible_to_module_form_fields.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xd9\xd4i\xe3W\x00\x00\x00W\x00\x00\x005\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xd4\x10\x00\x0020181202163130.fix-crm-module-form-primary-key.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\"\x96\xd6pj\x00\x00\x00j\x00\x00\x000\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x97\x11\x00\x0020181204123650.add-crm-content-json-field.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xb7\x93\xd4\xf6f\x00\x00\x00f\x00\x00\x004\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81h\x12\x00\x0020181204155326.add-crm-module-form-json-field.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(mA\xa8\x1e&\x02\x00\x00&\x02\x00\x00+\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x819\x13\x00\x0020181216214630.crm-content-to-record.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xcf\xc6g\xf6\xe4\x01\x00\x00\xe4\x01\x00\x00$\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xc1\x15\x00\x0020181217100000.add-charts-tbl.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xae \xfd2\x18\x00\x00\x00\x18\x00\x00\x00#\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x18\x00\x0020181224122301.rem-crm_field.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(+\xad\xb7\xed\xb8\x02\x00\x00\xb8\x02\x00\x00&\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81r\x18\x00\x0020190108100000.add-triggers-tbl.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x94#\xb9\x99-\x00\x00\x00-\x00\x00\x00/\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x87\x1b\x00\x0020190110175924.rem-crm-record-json-field.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x04]{\x1fo\x04\x00\x00o\x04\x00\x008\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x1a\x1c\x00\x0020190114072000.cleanup-record-tables-and-multival.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x8f\xecQt`\x01\x00\x00`\x01\x00\x00'\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xf8 \x00\x0020190121132408.record-updated-by.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(9S\x05%x\x01\x00\x00x\x01\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xb6\"\x00\x00migrations.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xc1h\xf1\xfb/\x00\x00\x00/\x00\x00\x00\x06\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xed\x81s$\x00\x00new.shUT\x05\x00\x01\x80Cm8PK\x05\x06\x00\x00\x00\x00\x12\x00\x12\x00\xbd\x06\x00\x00\xdf$\x00\x00\x00\x00"
}

func init() {
	fs.Register(Data())
}
