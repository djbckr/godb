create table [SYS].[_obj] (
    [_id] uuid,
    [_name] varchar(128)
);

create table [SYS].[_table] (
    [_id] uuid
)

create table [SYS].[_enum] (
    [_id] uuid
)

create table [SYS].[_enumValue] (
    [_id] uuid,
    [_enum_id] uuid,
    [_name] varchar(128),
    [_value] integer
)

create table [SYS].[_constraint] (

)



create view tables as
select *
  from [SYS].[_table] t
