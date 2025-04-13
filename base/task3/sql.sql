-- 向 students 表插入姓名为“张三”、年龄20、年级为“三年级”的记录：
INSERT INTO students (name, age, grade)
VALUES ('张三', 20, '三年级');

-- 检索所有年龄大于18岁的学生信息：
SELECT * FROM students WHERE age > 18;

-- 将姓名为“张三”的学生年级更新为“四年级”：
UPDATE students
SET grade = '四年级'
WHERE name = '张三';

-- 删除所有年龄小于15岁的学生记录：
DELETE FROM students
WHERE age < 15;

-- 转账事务
START TRANSACTION -- 开启事务
-- 检查账户A余额并扣款
UPDATE accounts
SET balance = balance - 100
WHERE id = 'A' AND balance >= 100;  -- 带条件的原子操作

-- 检查扣款是否成功
IF ROW_COUNT() = 0 THEN
    ROLLBACK;  -- 余额不足则回滚（网页7）
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '账户A余额不足';
END IF;

UPDATE accounts
SET balance = balance + 100
WHERE id = 'B';

-- 记录转账流水
INSERT INTO transactions (from_account_id, to_account_id, amount)
VALUES ('A', 'B', 100);

COMMIT;  -- 提交事务