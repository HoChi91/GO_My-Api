CREATE TABLE IF NOT EXISTS USERS(
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS AUTHORS (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  birth_date VARCHAR(10),
  description TEXT
);

CREATE TABLE IF NOT EXISTS BOOKS (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  author INT NOT NULL,
  publication_date VARCHAR(10),
  summary TEXT,
  stock INT DEFAULT 0,
  price DECIMAL(10, 2),
  FOREIGN KEY (author) REFERENCES AUTHORS(id)
);

CREATE TABLE IF NOT EXISTS PURCHASE_HISTORY (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user INT NOT NULL,
  book INT NOT NULL,
  quantity INT NOT NULL,
  total_price DECIMAL(10, 2) NOT NULL,
  payment_timestamp VARCHAR(50) NOT NULL,
  FOREIGN KEY (user) REFERENCES USERS(id),
  FOREIGN KEY (book) REFERENCES BOOKS(id)
);


INSERT INTO authors (id, name, birth_date, description) 
VALUES 
(1, 'Jane Austen', '1775-12-16', 'Auteur britannique célèbre pour ses romans comme "Orgueil et Préjugés".'),
(2, 'Mark Twain', '1835-11-30', 'Écrivain américain connu pour "Les Aventures de Tom Sawyer".'),
(3, 'Virginia Woolf', '1882-01-25', 'Pionnière du modernisme dans la littérature anglaise.'),
(4, 'Charles Dickens', '1812-02-07', 'Auteur anglais célèbre pour "Oliver Twist" et "David Copperfield".'),
(5, 'Mary Shelley', '1797-08-30', 'Écrivain britannique connue pour "Frankenstein".'),
(6, 'Franz Kafka', '1883-07-03', 'Écrivain tchèque auteur de "La Métamorphose".'),
(7, 'Fyodor Dostoevsky', '1821-11-11', 'Écrivain russe connu pour "Crime et Châtiment".'),
(8, 'George Orwell', '1903-06-25', 'Écrivain britannique connu pour "1984" et "La Ferme des animaux".'),
(9, 'Harper Lee', '1926-04-28', 'Auteur américaine célèbre pour "Ne tirez pas sur l oiseau moqueur".'),
(10, 'Gabriel Garcia Marquez', '1927-03-06', 'Auteur colombien célèbre pour "Cent ans de solitude".');


INSERT INTO BOOKS (title, author, publication_date, summary, stock, price) VALUES
-- Livres de Jane Austen
('Pride and Prejudice', 1, '1813-01-28', 'Un roman d’amour et de société.', 10, 12.99),
('Sense and Sensibility', 1, '1811-10-30', 'Un roman sur les émotions et la raison.', 8, 10.99),
('Emma', 1, '1815-12-25', 'Un récit captivant sur l’aristocratie rurale.', 7, 15.49),
('Mansfield Park', 1, '1814-07-09', 'Une exploration des classes sociales.', 9, 11.99),
('Northanger Abbey', 1, '1817-12-01', 'Une parodie de romans gothiques.', 6, 9.99),

-- Livres de Mark Twain
('The Adventures of Tom Sawyer', 2, '1876-06-17', 'Un roman sur l’enfance et l’aventure.', 12, 13.99),
('The Adventures of Huckleberry Finn', 2, '1885-12-10', 'Un classique de la littérature américaine.', 11, 14.99),
('The Prince and the Pauper', 2, '1881-11-19', 'Une histoire fascinante d’échange d’identités.', 8, 10.99),
('A Connecticut Yankee in King Arthur’s Court', 2, '1889-04-21', 'Une satire hilarante sur le temps.', 5, 9.99),
('Pudd’nhead Wilson', 2, '1894-11-28', 'Une analyse incisive des préjugés.', 7, 11.49),

-- Livres de Virginia Woolf
('Mrs. Dalloway', 3, '1925-05-14', 'Une journée dans la vie de Clarissa Dalloway.', 10, 14.99),
('To the Lighthouse', 3, '1927-05-05', 'Un roman introspectif sur la famille.', 6, 12.99),
('Orlando', 3, '1928-10-11', 'Une biographie fictive traversant les siècles.', 9, 13.99),
('The Waves', 3, '1931-10-08', 'Un récit expérimental sur la conscience.', 8, 15.49),
('A Room of One’s Own', 3, '1929-09-24', 'Un essai féministe révolutionnaire.', 7, 11.99),

-- Livres de Charles Dickens
('Oliver Twist', 4, '1837-02-01', 'Les aventures d’un orphelin à Londres.', 15, 10.99),
('David Copperfield', 4, '1850-05-01', 'Un roman d’apprentissage classique.', 12, 13.49),
('Great Expectations', 4, '1861-08-01', 'L’histoire d’un jeune homme ambitieux.', 11, 14.99),
('A Tale of Two Cities', 4, '1859-11-01', 'Un récit sur la Révolution française.', 10, 12.49),
('Bleak House', 4, '1853-03-01', 'Une satire du système judiciaire anglais.', 8, 15.99),

-- Livres de Mary Shelley
('Frankenstein', 5, '1818-01-01', 'Le classique du roman gothique et d’horreur.', 10, 12.99),
('The Last Man', 5, '1826-01-01', 'Un récit post-apocalyptique.', 6, 11.99),
('Valperga', 5, '1823-01-01', 'Une histoire d’amour dans l’Italie médiévale.', 7, 13.49),
('Lodore', 5, '1835-01-01', 'Une exploration de la maternité et de l’autonomie.', 8, 10.99),
('Matilda', 5, '1959-01-01', 'Un récit inachevé sur le chagrin et la rédemption.', 5, 9.99),

-- Livres de Franz Kafka
('The Metamorphosis', 6, '1915-10-25', 'Une métaphore de l’aliénation.', 9, 12.49),
('The Trial', 6, '1925-04-26', 'Un récit kafkaïen sur l’absurdité bureaucratique.', 8, 13.99),
('The Castle', 6, '1926-08-10', 'Un homme aux prises avec un système oppressant.', 7, 14.49),
('In the Penal Colony', 6, '1919-06-22', 'Une réflexion sur la justice et la cruauté.', 6, 10.99),
('Amerika', 6, '1927-09-21', 'Un roman inachevé sur l’immigration.', 8, 11.99),

-- Livres de Fyodor Dostoevsky
('Crime and Punishment', 7, '1866-01-01', 'Un chef-d’œuvre psychologique.', 12, 14.99),
('The Brothers Karamazov', 7, '1880-01-01', 'Un roman sur la foi et la famille.', 10, 15.49),
('Notes from Underground', 7, '1864-01-01', 'Une exploration du libre arbitre.', 9, 11.99),
('Demons', 7, '1872-01-01', 'Un roman politique et philosophique.', 8, 13.99),
('The Idiot', 7, '1869-01-01', 'Un homme bon dans un monde corrompu.', 7, 12.99),

-- Livres de George Orwell
('1984', 8, '1949-06-08', 'Une dystopie sur la surveillance et le totalitarisme.', 15, 16.99),
('Animal Farm', 8, '1945-08-17', 'Une allégorie politique sur le communisme.', 14, 12.49),
('Homage to Catalonia', 8, '1938-01-01', 'Un récit autobiographique sur la guerre civile espagnole.', 10, 14.99),
('Down and Out in Paris and London', 8, '1933-01-01', 'Une immersion dans la pauvreté.', 9, 10.99),
('Coming Up for Air', 8, '1939-01-01', 'Une réflexion sur le changement et la mémoire.', 8, 11.49),

-- Livres de Harper Lee
('To Kill a Mockingbird', 9, '1960-07-11', 'Un roman sur l’injustice raciale.', 12, 15.99),
('Go Set a Watchman', 9, '2015-07-14', 'Une suite du célèbre roman.', 10, 13.49),
('Atticus Finch: The Biography', 9, '1965-01-01', 'Un regard fictif sur le protagoniste.', 8, 11.99),
('Maycomb Stories', 9, '1962-01-01', 'Un recueil de nouvelles sur le Sud profond.', 7, 10.99),
('Mockingbird Revisited', 9, '1970-01-01', 'Un hommage au roman original.', 6, 12.49),

-- Livres de Gabriel Garcia Marquez
('One Hundred Years of Solitude', 10, '1967-06-05', 'Une épopée magique et historique.', 15, 17.99),
('Love in the Time of Cholera', 10, '1985-09-15', 'Une histoire d’amour intemporelle.', 13, 15.49),
('Chronicle of a Death Foretold', 10, '1981-01-01', 'Un roman sur un meurtre annoncé.', 12, 13.99),
('The General in His Labyrinth', 10, '1989-01-01', 'Un portrait fictif de Simon Bolivar.', 10, 14.49),
('Of Love and Other Demons', 10, '1994-01-01', 'Une histoire sur l’amour et les préjugés.', 9, 12.99);
