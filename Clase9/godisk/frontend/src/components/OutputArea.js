import React from 'react';

const OutputArea = ({ output, onClear }) => {
  return (
    <div className="card">
      <div className="card-header d-flex justify-content-between align-items-center">
        <h5 className="card-title mb-0">
          <i className="bi bi-list-ul me-2"></i>
          Área de Salida de Comandos
        </h5>
        <button
          type="button"
          className="btn btn-outline-secondary btn-sm"
          onClick={onClear}
        >
          <i className="bi bi-trash me-1"></i>
          Limpiar
        </button>
      </div>
      <div className="card-body">
        <pre
          className="bg-dark text-light p-3 rounded"
          style={{ 
            height: '500px', 
            overflowY: 'auto',
            fontSize: '0.875rem',
            whiteSpace: 'pre-wrap'
          }}
        >
          {output || 'Ejecute algunos comandos para ver la salida aquí...'}
        </pre>
      </div>
    </div>
  );
};

export default OutputArea;