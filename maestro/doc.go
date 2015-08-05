/*

apihub-client is an open source command line solution for publishing APIs on ApiHub Servers.

Usage:
  % apihub command [command options] [arguments...]

The currently available commands are:

  login               Login in with your ApiHub credentials.
  logout              Clear local credentials.

  service-add         Create a new service.
  service-remove      Remove an existing service.

  target-add          Add a new target in the list of targets.
  target-list         Adds a new target in the list of targets.
  target-remove       Remove a target from the list of targets.
  target-set          Set a target as default.

  team-create         Create a team.
  team-info           Return team info and lists of members and services.
  team-list           Return a list of all teams.
  team-remove         Delete a team.
  team-user-add       Add a user to a team.
  team-user-remove    Remove a user from a team.

  user-create         Create a user account.
  user-remove         Delete a user account.

  GLOBAL OPTIONS:
   --version, -v   print the version

Use "apihub <command> --help, or -h" for more information about a command.


Authentication on ApiHub server

Usage:

  % apihub login <email>

The email address and password are used by the client to obtain an API token. This token is used for authentication in the following api requests, for the elected target.
The apihub-client stores the token in the standard Unix file ~/.apihub_token and the content looks like:
  % cat ~/.apihub_token
  % Token iBnD0Epiz4pX1zNDYGLhUpjnF33mvElvfIGTzSFuuVc=


Logout from ApiHub server

Usage:

  % apihub logout

Clear local token (the file file ~/.apihub_token).


Create a new service

USAGE:

  % apihub service-add --team <team> --subdomain <subdomain> --endpoint <api_endpoint>

OPTIONS

  --description,   -desc  Service description
  --disabled,      -dis   Disable the service
  --documentation, -doc   Url with the documentation
  --endpoint,      -e     Url where the service is running

  --keyless,       -k     Indicate if it's allowed to make requests without authentication.

  --subdomain,     -s     The subdomain will be used by the proxy: http://ratings.apihubserver.org. (where `ratings` is the chosen subdomain).

  --team,           -t    Team responsible for the service
  --timeout         Timeout Default: 0 (Do nothing. Wait the api server to return timeout.)


Remove an existing service

USAGE:

  % apihub service-remove --subdomain <subdomain>

OPTIONS:
   --subdomain, -s  Subdomain

This action cannot be undone. Once a service is deleted, it's needed to add and configure it again.


Manage ApiHub server endpoints

USAGE:

  % apihub target-list
  % apihub target-add <label> <endpoint>
  % apihub target-set <label>
  % apihub target-remove <label>

Target is the ApiHub server endpoint. It's possible to have multiple instances runnning and still use the same apihub-client. You just need to add a new target and mark it as default, by using the commands `target-add` and `target-set` respectively.
Requests operations will be directed to the elected target. It's possible to check the current target by using the command `target-list`.

File format

The file contains a list of all targets and a flag indicating what is the current:

  current: home
  options:
    home: http://api.example.com
    apihub: http://github.com/apihub


Create a new team

USAGE:

  % apihub team-create --name <name>

OPTIONS:
   --name, -n   Name of the team

By creating a new team, the current user is added to it automatically as owner. It's created an `alias`(slug) based on the name of the team.
You should use the alias when interacting with other team commands.


Return team info and lists of members and services

USAGE:

  % apihub team-info --alias <alias>

OPTIONS:
   --alias, -a  Team alias

The `alias` is a slug version of the team name. If you do not know it, run `apihub team-list` for more details.
  Name: ApiHub
  Alias: apihub
  Owner: bob@sample.org

  +----------------+
  |  TEAM MEMBERS  |
  +----------------+
  | bob@sample.org |
  +----------------+


Return a list of all teams

USAGE:

  % apihub team-list

Return a list containing all the teams you belong to and the owner for each of them.
  +-----------+-----------+----------------+
  | TEAM NAME |   ALIAS   |     OWNER      |
  +-----------+-----------+----------------+
  | ApiHub | apihub | bob@sample.org |
  +-----------+-----------+----------------+


Delete a team

USAGE:

  % apihub team-remove --alias <alias>

OPTIONS:
   --alias, -a  Team alias

The `alias` is a slug version of the team name. If you do not know it, run `apihub team-list` for more details.


Add a user to a team

USAGE:

  % apihub team-user-add --team <team-alias> --email <email>

OPTIONS:
   --team, -t   Name of the team
   --email, -e  User's email

You need to be a member of the team to add another user.


Remove a user from a team

USAGE:

  % apihub team-user-add --team <team-alias> --email <email>

OPTIONS:
   --team,  -t  Name of the team
   --email, -e  User's email

You need to be a member of the team to remove a user from a team.


Create a user account

USAGE:

  % apihub  user-create --name <name> --email <email>

OPTIONS:
   --name, -n   The user's real life name
   --email, -e  User's email

Creates a new account on the ApiHub server. It's important to notice that the account is created only on the current target.
If you use multiple server instances, you need to create one account for each of them.


Delete a user account

USAGE:

  % apihub  user-remove

Delete current logged in account from ApiHub server and deletes the file `~/.apihub_token`, which contain the token.


Return the current apihub-client version

USAGE:

  % apihub  --version, -v

Return the current version of the client. This version will be sent via header `ApiHubClient-Version` to all API requests.

*/
package apihub
