import axios from 'axios';

const API_BASE_URL = 'http://localhost:3001'; // ajustar al puerto correcto del backend

export const executeCommands = async (commands) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/execute`, {
      commands: commands
    });
    return response.data.output;
  } catch (error) {
    throw new Error('Error al comunicarse con el servidor: ' + error.message);
  }
};