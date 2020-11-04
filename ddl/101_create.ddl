/* ------------------------------------- マスタ群 ------------------------------------- */

/* 【マスタ】カテゴリー */
CREATE TABLE accountBook.m_categorys (
    id TINYINT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    name VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);

/* 【マスタ】消費税 */
/*
    CREATE TABLE accountBook.m_tax(
    id TINYINT UNSIGNED AUTO_INCREMENT NOT NULL,
    tax FLOAT UNSIGNED NOT NULL UNIQUE,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);
*/

/* 【マスタ】ユーザ */
CREATE TABLE accountBook.m_users(
    id INT(5) UNSIGNED ZEROFILL NOT NULL,
    name VARCHAR(10) NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);

/* 【マスタ】購入品ステータス */
CREATE TABLE accountBook.m_item_status (
    id TINYINT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    status VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);

/* 【マスタ】使用タイミング */
CREATE TABLE accountBook.m_used_timing (
    id TINYINT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    timing VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);


/* ------------------------------------- テーブル群 ------------------------------------- */

/* 【テーブル】月予算 */
CREATE TABLE accountBook.t_budgets(
    id TINYINT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    month VARCHAR(7) NOT NULL UNIQUE,
    budget SMALLINT UNSIGNED,
    created_at TIMESTAMP NOT NULL,
    created_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY fk_created_user_id(created_user_id) references m_users(id),
    FOREIGN KEY fk_updated_user_id(updated_user_id) references m_users(id)
);


/* 【テーブル】購入品 */
CREATE TABLE accountBook.t_purchased_items(
    id INT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    month VARCHAR(7) NOT NULL,
    purchased_on DATE NOT NULL,
    shop VARCHAR(30),
    category_id TINYINT UNSIGNED ZEROFILL,
    item_name VARCHAR(30) NOT NULL,
    quantity SMALLINT UNSIGNED NOT NULL,
    price_including_tax SMALLINT UNSIGNED NOT NULL,
    item_status_id TINYINT UNSIGNED ZEROFILL NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY fk_category_id(category_id) references m_categorys(id),
    FOREIGN KEY fk_item_status_id(item_status_id) references m_item_status(id),
    FOREIGN KEY fk_created_user_id(created_user_id) references m_users(id),
    FOREIGN KEY fk_updated_user_id(updated_user_id) references m_users(id)
);

/* 【テーブル】使用品 */
CREATE TABLE accountBook.t_used(
    id INT UNSIGNED ZEROFILL AUTO_INCREMENT NOT NULL,
    used_on DATE NOT NULL,
    used_timing_id TINYINT UNSIGNED ZEROFILL,
    purchase_item_id INT UNSIGNED ZEROFILL NOT NULL,
    quantity SMALLINT UNSIGNED NOT NULL,
    menu VARCHAR(30),
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_user_id INT(5) ZEROFILL UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY fk_used_timing_id(used_timing_id) references m_used_timing(id),
    FOREIGN KEY fk_purchase_item_id(purchase_item_id) references t_purchased_items(id),
    FOREIGN KEY fk_created_user_id(created_user_id) references m_users(id),
    FOREIGN KEY fk_updated_user_id(updated_user_id) references m_users(id)
);

