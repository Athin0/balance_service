CREATE DATABASE  MyDB;
CREATE TABLE IF NOT EXISTS Hall
(
    HallID SERIAL NOT NULL,
    number_seats INT CHECK (number_seats >= 0) DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    PRIMARY KEY (HallID)
    );

CREATE TABLE IF NOT EXISTS Seat
(
    SeatID SERIAL NOT NULL,
    row_number INT NOT NULL CHECK(row_number>0),
    seat_number INT NOT NULL CHECK(seat_number>0),
    vip BOOLEAN default false,
    HallID INT NOT NULL,
    UNIQUE (seat_number, row_number,HallID),
    PRIMARY KEY (SeatID),
    FOREIGN KEY (HallID) REFERENCES Hall(HallID)
    );

CREATE TABLE IF NOT EXISTS Person
(
    PersonID SERIAL NOT NULL,
    password VARCHAR(200) CHECK(char_length(password)>=6)NOT NULL,
    login VARCHAR(50)  CHECK (char_length(login) >= 6) UNIQUE NOT NULL,
    position varchar(64) NOT NULL,
    PRIMARY KEY (PersonID)
    );


CREATE TABLE IF NOT EXISTS Point
(
    PointID SERIAL NOT NULL,
    name VARCHAR(64) NOT NULL,
    address VARCHAR(124) NOT NULL,
    PRIMARY KEY (PointID)
    );


CREATE TABLE IF NOT EXISTS Country
(
    id_country SERIAL NOT NULL ,
    country_name varchar(64) NOT NULL UNIQUE ,
    PRIMARY KEY (id_country)
    );

CREATE TABLE IF NOT EXISTS Film
(
    id_film serial NOT NULL,
    tittle VARCHAR(64) NOT NULL,
    genre varchar(64) NOT NULL ,
    main_roles text ,
    description text NOT NULL ,
    producer varchar(256),
    duration INTERVAL HOUR TO MINUTE NOT NULL,    /*todo  какой тип*/
    age INT CHECK ( age>=0 and age<40) NOT NULL,
    PRIMARY KEY (id_film)
    );
CREATE TABLE IF NOT EXISTS Country_film
(
    id_country INT NOT NULL,
    id_film INT NOT NULL,
    PRIMARY KEY (id_country, id_film),
    FOREIGN KEY (id_country) REFERENCES Country(id_country),
    FOREIGN KEY (id_film) REFERENCES Film(id_film)
);


CREATE TABLE IF NOT EXISTS Seance
(
    SeanceID serial NOT NULL,
    date timestamp NOT NULL,       /*day and time*/
    format varchar(10)  NOT NULL ,         /*CHECK (format in ('2D', '3D', 'IMAX')) enum*/
    price_simple INT check ( price_simple>=0 ) NOT NULL,
    price_vip  INT check (price_vip>=0) NOT NULL,
    id_film INT NOT NULL,
    HallID INT NOT NULL,
    PointID INT NOT NULL,
    PRIMARY KEY (SeanceID),
    FOREIGN KEY (id_film) REFERENCES Film(id_film),
    FOREIGN KEY (HallID) REFERENCES Hall(HallID),
    FOREIGN KEY (PointID) REFERENCES Point(PointID)
    );

CREATE TABLE IF NOT EXISTS Ticket
(
    TicketId serial NOT NULL,
    booked boolean NOT NULL DEFAULT false,
    paid boolean NOT NULL DEFAULT false,
    SeatID INT NOT NULL,
    PersonID INT NOT NULL,
    SeanceID INT NOT NULL,
    PRIMARY KEY (TicketId),
    FOREIGN KEY (SeatID) REFERENCES Seat(SeatID),
    FOREIGN KEY (PersonID) REFERENCES Person(PersonID),
    FOREIGN KEY (SeanceID) REFERENCES Seance(SeanceID)
    );


