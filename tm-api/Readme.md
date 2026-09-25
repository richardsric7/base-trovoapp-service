# Project Repository Documentation

Welcome to our project repository. This README provides instructions on setting up the project environment and importing the necessary SQL backups using pgAdmin.

## Setting Up Environment Variables

Before starting, you'll need to configure your environment variables.

1. **Using the `.env-sample` File:**
    - Locate the `.env-sample` file in the repository.
    - This file contains template environment variables required for the project.
    - Copy the contents of `.env-sample` into a new file and name it `.env`.
    - Replace the placeholder values in `.env` with your actual configuration values. This includes database credentials and other necessary configurations.

## Importing SQL Backups in pgAdmin

To set up your database, you'll need to import the SQL backup files.

### Backup Files

- The SQL backup files used for this project are:
    - `trovop2p_wip.sql`
    - `trovowalletwip.sql`
- The two files can be found in the root of this project

### Importing Steps

1. **Open pgAdmin:**
    - Launch pgAdmin and connect to your database server.

2. **Create two New Databases for this project: (Optional)**
    - One is for `trovowallet`, the other is for `trovop2p` the two dbs are needed for the project to work properly.
    - Right-click on 'Databases', choose 'Create', then 'Database', and follow the prompts.

3. **Importing the SQL File:**
    - Download the two databases required to run this project (`trovop2p_wip.sql` or `trovowalletwip.sql`)
    - Right-click on the database you want to import the backup into.
    - Navigate to 'Restore' or 'Import', depending on your pgAdmin version.
    - In the dialog box, click on the 'Upload' button and choose either `trovop2p_wip.sql` or `trovowalletwip.sql`.
    - Follow the on-screen instructions to complete the import process.

4. **Update Database Credentials:**
    - Ensure that the database credentials in your `.env` file match the database you've imported the SQL files into.

### Alternative Tutorial

- If you encounter any issues or need a visual guide, you can watch this YouTube tutorial: [How to Import a SQL Backup in pgAdmin](https://www.youtube.com/watch?v=4HJwrXclC3A).

## Final Steps

After setting up your environment variables and importing the SQL backups, you should be ready to start using the project. Remember to keep your `.env` file secure and not to upload it to any public repositories.

## Support

For additional help or information, feel free to contact Toluwase.

 doctl registry login&& docker build -t registry.digitalocean.com/service-images/trovo-wallet-api:dsb-dev . &&docker push registry.digitalocean.com/service-images/trovo-wallet-api:dsb-dev