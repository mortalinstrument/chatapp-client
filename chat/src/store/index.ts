import { createStore } from 'vuex'

export type Message = {
  Message: string,
  Timestamp: string,
  From: User
}

export type minifiedMessage = {
  Message: string,
  Timestamp: string,
  From: string
}

export type User = {
  Name: string,
  IP: string,
  LastLogin: string,
  Active: boolean
  Messages: minifiedMessage[]
}

const socket = new WebSocket("ws://localhost:7777/c")

export interface partnersState {
  partners: Array<User>
}

export default createStore({
  state: (): partnersState => ({
    partners: Array<User>()
  }),
  getters: {
    getPartners: state => {
      return state.partners
    }
  },
  mutations: {
    addMessageFromWebsocket(state, payload:Message){
      const index = state.partners.findIndex((partner) => partner.IP == payload.From.IP )
      if(index == -1){
        //add partner and then messages
        console.log(payload.From)
        state.partners.push(payload.From)
        const index = state.partners.findIndex((partner) => partner.IP == payload.From.IP )
        state.partners[index].Messages = new Array<minifiedMessage>
        const messageToSave : minifiedMessage = {
          Message: payload.Message,
          Timestamp: payload.Timestamp,
          From: payload.From.Name
        }
        state.partners[index].Messages.unshift(messageToSave)
      } else {
        //just add messages
        const index = state.partners.findIndex((partner) => partner.IP == payload.From.IP )
        const messageToSave : minifiedMessage = {
          Message: payload.Message,
          Timestamp: payload.Timestamp,
          From: payload.From.Name
        }
        state.partners[index].Messages.unshift(messageToSave)
      }
    },
    addMessageFromFrontend(state, payload: minifiedMessage){
      const index = state.partners.findIndex((partner) => partner.Name == payload.From )
      console.log(payload)
      state.partners[index].Messages.unshift(payload)
    }
  },
  actions: {
    async establishWebsocketConnection(context){
      socket.onopen = () => {
        socket.send("Client is here!")
      }

      socket.onclose = event => {
        console.log("Socket closed: " + event.reason)
      }

      socket.onerror = error => {
        console.log("Socket error: " + error)
      }

      socket.onmessage = await function (evt){
        const jsonObject = JSON.parse(evt.data)
        console.log(jsonObject)
        context.commit('addMessageFromWebsocket', jsonObject)
      }
    },
    sendMessage(context, message: minifiedMessage){
      socket.send(JSON.stringify(message))
      console.log(JSON.stringify(message))
      context.commit('addMessageFromFrontend', message)
    }
  },
  modules: {
  }
})
