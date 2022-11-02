cp .env.example .env
sed -i -e "s/<your tag here>/$1/g" .env