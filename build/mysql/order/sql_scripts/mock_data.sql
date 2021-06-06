CREATE DEFINER=`root`@`%` FUNCTION `mock_data`() RETURNS int
    DETERMINISTIC
begin
    declare num int default 1000000;
    declare i int default 0;
    while i<num do
            insert into users(`name`, `email`, `passwd`)
            values (concat('user', i), concat(i, '@gmail.com'), sha2(concat('123456', 'salt'), 256));
            set i=i+1;
        end while;
    return i;
end;

select mock_data();