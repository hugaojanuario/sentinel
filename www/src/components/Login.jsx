import { useState } from 'react'
import { login, register } from '../api'

export default function Login({ onLogin }) {
  const [mode, setMode] = useState('login')
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState(null)
  const [success, setSuccess] = useState(null)
  const [loading, setLoading] = useState(false)

  function switchMode(next) {
    setMode(next)
    setError(null)
    setSuccess(null)
    setName('')
    setEmail('')
    setPassword('')
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError(null)
    setSuccess(null)
    setLoading(true)
    try {
      if (mode === 'login') {
        const token = await login(email, password)
        onLogin(token)
      } else {
        await register(name, email, password)
        setSuccess('Conta criada! Faça login.')
        switchMode('login')
      }
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-wrap">
      <form className="login-form" onSubmit={handleSubmit}>
        <div className="login-logo">
          <img
            src="https://github.com/user-attachments/assets/35f3d090-1e6a-40c1-9621-894ff89004ea"
            alt="Sentinel"
            width="40"
            height="40"
          />
          <span className="logo-text">Sentinel</span>
        </div>
        {mode === 'register' && (
          <div className="login-field">
            <label>Nome</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Seu nome"
              required
              autoFocus
            />
          </div>
        )}
        <div className="login-field">
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="seu@email.com"
            required
            autoFocus={mode === 'login'}
          />
        </div>
        <div className="login-field">
          <label>Senha</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            required
          />
        </div>
        {error && <p className="login-error">{error}</p>}
        {success && <p className="login-success">{success}</p>}
        <button type="submit" className="login-btn" disabled={loading}>
          {loading
            ? mode === 'login' ? 'Entrando...' : 'Criando...'
            : mode === 'login' ? 'Entrar' : 'Criar conta'}
        </button>
        <p className="login-toggle">
          {mode === 'login' ? (
            <>Não tem conta?{' '}
              <button type="button" onClick={() => switchMode('register')}>Registrar</button>
            </>
          ) : (
            <>Já tem conta?{' '}
              <button type="button" onClick={() => switchMode('login')}>Entrar</button>
            </>
          )}
        </p>
      </form>
    </div>
  )
}
