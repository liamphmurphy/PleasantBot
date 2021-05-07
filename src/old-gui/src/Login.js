import React from 'react';

import 'bootstrap/dist/css/bootstrap.min.css';
import './css/style.css';

function Login() {
  return (
    <div>
        <a href="https://id.twitch.tv/oauth2/authorize?client_id=jwp1l9234hl647c13cyvbue5j6vequ&redirect_uri=http://localhost:3000&response_type=token&scope=chat:read%20chat:edit%20user:edit%20moderation:read">Login to Twitch</a>
    </div>
  );
}

export default Login;
