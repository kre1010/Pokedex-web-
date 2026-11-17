GET/api/v1/pokemon il y a 1025 pokemon les nom français sont ranger de cette manière  "name": {
            "fr": "Pêchaminus",
            "en": "Pecharunt",
            "jp": "モモワロウ"
        } 
le premier pokemon n'est pas un vrai pokemon c'est un pokemon bug

Get/api/v1/pokemon/25 la fiche de pikachu s'affiche bien {
    "pokedex_id": 25,
    "generation": 1,
    "category": "Pokémon Souris",
    "name": {
        "fr": "Pikachu",
        "en": "Pikachu",
        "jp": "ピカチュウ"
    },

GET/api/v1/pokemon/9999 quand un pokemon n'existe Postman affiche {
    "status": 404,
    "message": "Impossible de trouver le Pokémon dans la base de données."
}