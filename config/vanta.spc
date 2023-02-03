connection "vanta" {
  plugin = "vanta"

  # A personal API token to access Vanta API
  # This is only required while querying `vanta_evidence` table. 
  # To generate an API token, refer: https://developer.vanta.com/docs/quick-start#1-make-an-api-token
  # api_token = "97GtVsdAPwowRToaWDtgZtILdXI_agszONwajQslZ1o"

  # Session ID of your current vanta session
  # Set the value of `connect.sid` cookie from a logged in Vanta browser session
  # Required to access tables that are using the https://app.vanta.com/graphql endpoint
  # session_id = "s:3nZSteamPipe1fSu4iNV_1TB5UTesTToGK.zVANtaplugintest+GVxPvQffhnFY3skWlfkceZxXKSCjc"
}