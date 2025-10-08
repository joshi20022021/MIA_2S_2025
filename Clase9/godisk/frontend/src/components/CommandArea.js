import React, { useState } from 'react';

const CommandArea = ({ onExecute, isExecuting }) => {
  const [commands, setCommands] = useState('');

  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setCommands(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  const handleExecute = () => {
    if (commands.trim()) {
      onExecute(commands);
    }
  };


  return (
    <div className="card">
      <div className="card-header">
        <h5 className="card-title mb-0">
          <i className="bi bi-terminal me-2"></i>
          Área de Entrada de Comandos
        </h5>
      </div>
      <div className="card-body">
        <div className="mb-3">
          <label htmlFor="fileInput" className="form-label">
            Cargar archivo de comandos:
          </label>
          <input
            type="file"
            className="form-control"
            id="fileInput"
            accept=".txt,.sh"
            onChange={handleFileUpload}
          />
        </div>
        
        <div className="mb-3">
          <label htmlFor="commandsTextarea" className="form-label">
            Comandos:
          </label>
          <textarea
            className="form-control font-monospace"
            id="commandsTextarea"
            rows="15"
            value={commands}
            onChange={(e) => setCommands(e.target.value)}
            placeholder="Ingrese los comandos aquí..."
          />
        </div>

        <div className="d-grid gap-2 d-md-flex justify-content-md-end">
          <button
            type="button"
            className="btn btn-success"
            onClick={handleExecute}
            disabled={isExecuting || !commands.trim()}
          >
            {isExecuting ? (
              <>
                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                Ejecutando...
              </>
            ) : (
              <>
                <i className="bi bi-play-fill me-1"></i>
                Ejecutar
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default CommandArea;