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

export type outgoingMessage = {
  Message: string,
  Timestamp: string,
  ToIP: string,
  ToName: string
}

export type User = {
  Name: string,
  IP: string,
  LastLogin: string,
  Active: boolean
  Messages: minifiedMessage[]
}

const messageSocket = new WebSocket("ws://localhost:7777/c")
const userSocket =  new WebSocket("ws://localhost:7777/cu")


export interface partnersState {
  partners: Array<User>
  possiblePartners: Array<User>
  myself: User
}

export default createStore({
  state: (): partnersState => ({
    partners: Array<User>(),
    possiblePartners: Array<User>(),
    myself: {} as User
  }),
  getters: {
    getPartners: state => {
      return state.partners
    },
    getPossiblePartners: state => {
      return state.possiblePartners
    },
    getUserInfo: state => {
      return state.myself
    }
  },
  mutations: {
    newChat(state, newPartnerToBe: User){
      const index = state.partners.findIndex((partner) => partner.IP == newPartnerToBe.IP )
      if(index == -1){
        newPartnerToBe.Messages = new Array<minifiedMessage>
        state.partners.unshift(newPartnerToBe)
      }
    },
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
    addMessageFromFrontend(state, payload: outgoingMessage){
      const index = state.partners.findIndex((partner) => partner.Name == payload.ToName )
      console.log(index, payload)
      const msgToSave: minifiedMessage = {
        Message: payload.Message,
        Timestamp: payload.Timestamp,
        From: state.myself.Name 
      }
      state.partners[index]?.Messages.unshift(msgToSave)
    },
    addPossiblePartnerFromWebsocket(state, payload: User){
      state.possiblePartners.unshift(payload)
    },
    removePossiblePartnerFromWebsocket(state, payload: string){
      state.possiblePartners = state.possiblePartners.filter(function(value, index, arr){ 
        return value.IP != payload;
    });
    },
    setUserData(state, payload: User) {
      state.myself = payload
    },
    setPossiblePartners(state, payload: User[]){
      state.possiblePartners = payload
    }
  },
  actions: {
    async establishMessageWebsocketConnection(context){
      messageSocket.onclose = event => {
        console.log("Message Socket closed: " + event.reason)
      }

      messageSocket.onerror = error => {
        console.log("Message Socket error: " + error)
      }

      messageSocket.onmessage = await function (evt){
        const jsonObject = JSON.parse(evt.data)
        console.log("Message Websocket:" + jsonObject)
        context.commit('addMessageFromWebsocket', jsonObject)
      }
    },
    async establishUserWebsocketConnection(context){
      //first get all users, then init websocket
      const allUsersRequest = await fetch("http://localhost:7777/whothere")
      let allUsersData = await allUsersRequest.json()

      if (!allUsersRequest.ok) {
        const error = new Error(
          "getting all users failed"
        )
        throw error
      }

      if (allUsersData == null){
        allUsersData = Array<User>()
      }
      context.commit('setPossiblePartners', allUsersData)

      //setup event listeners for constant communication
      userSocket.onclose = event => {
        console.log("User Socket closed: " + event.reason)
      }

      userSocket.onerror = error => {
        console.log("User Socket error: " + error)
      }

      userSocket.onmessage = await function(evt){
        const evtString = String(evt.data)
        if(evtString.search("remove") != -1){
          const ip = evtString.split(":")[1]
          console.log("User Websocket: REMOVE: " + ip)
          context.commit('removePossiblePartnerFromWebsocket', ip)

        } else {
          const jsonObject = JSON.parse(evt.data)
          console.log("User Websocket: ADD: " + jsonObject)
          context.commit('addPossiblePartnerFromWebsocket', jsonObject)
        }
      }
    },
    sendMessage(context, message: outgoingMessage){
      messageSocket.send(JSON.stringify(message))
      context.commit('addMessageFromFrontend', message)
    },
    async requestUserInfo(context){
      const responseUser = await fetch("http://localhost:7777/whoami");
      const responseUserData = await responseUser.json();
  
      if (!responseUser.ok) {
        const error = new Error(
          "getting user info failed"
        )
        throw error
      }
      console.log(responseUserData)

      context.commit('setUserData', responseUserData)
    }
  },
  modules: {
  }
})