import { writable } from 'svelte/store'

let storedUser = null

// 确保在浏览器端访问localStorage
if (typeof window !== 'undefined') {
  storedUser = localStorage.getItem('user')
}

const userInformation = writable(storedUser ? JSON.parse(storedUser) : null)

// 监听userInformation的变化并同步到localStorage
userInformation.subscribe(($userInformation) => {
  if (typeof window !== 'undefined') {
    if ($userInformation) {
      localStorage.setItem('user', JSON.stringify($userInformation))
    } else {
      localStorage.removeItem('user')
    }
  }
})

const setUserInformation = (userData) => {
  userInformation.set(userData)
}

const removeUserInformation = () => {
  userInformation.set(null)
}

export const information = {
  userInformation,
  setUserInformation,
  removeUserInformation
}