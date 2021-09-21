create database homestead;
create user homestead with encrypted password 'secret';
grant all privileges on database homestead to homestead;
