CREATE TABLE login_history (
                               id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                               user_id BIGINT UNSIGNED NOT NULL,

    -- 登录信息
                               login_time DATETIME DEFAULT CURRENT_TIMESTAMP,
                               login_ip VARCHAR(45) NOT NULL COMMENT '支持IPv6',
                               user_agent TEXT NULL COMMENT '浏览器/设备信息',
                               login_method VARCHAR(20) COMMENT 'password, oauth, sso等',

    -- 登录结果
                               success BOOLEAN DEFAULT TRUE,
                               failure_reason VARCHAR(100) NULL COMMENT '失败原因',

    -- 位置信息（可选）
                               country VARCHAR(50) NULL,
                               region VARCHAR(50) NULL,
                               city VARCHAR(50) NULL,
                               latitude DECIMAL(10, 8) NULL,
                               longitude DECIMAL(11, 8) NULL,

    -- 设备指纹（安全相关）
                               device_id VARCHAR(100) NULL COMMENT '设备唯一标识',
                               session_id VARCHAR(100) NULL,

    -- 索引
                               INDEX idx_user_id (user_id),
                               INDEX idx_login_time (login_time),
                               INDEX idx_login_ip (login_ip),
                               INDEX idx_success (success),
                               FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;