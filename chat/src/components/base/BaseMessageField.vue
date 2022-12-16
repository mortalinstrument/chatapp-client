<template>
    <form @submit.prevent="send(partner)">
        <p v-if="!messageValid">Bitte geben sie eine Nachricht ein!</p>
        <ion-label position="floating">Neue Nachricht verfassen</ion-label>
        <ion-input placeholder="Hier tippen..." v-model="messageToSend"></ion-input>
        <ion-button type="submit">Send</ion-button>
    </form>
</template>

<script lang="ts">
    import { IonLabel, IonInput, IonButton } from '@ionic/vue'
    import { defineComponent } from 'vue';
    import { mapActions } from 'vuex';
    import { User, outgoingMessage } from '../../store/index'

    export default defineComponent ({
        components: { IonLabel, IonInput, IonButton },
        data() {
        return {
          messageToSend: "",
          messageValid: true
        }
      },
      props: ['partner'],
      methods:{
        send(partner: User){
          this.check()
          if(this.messageValid == false){
            return
          }
          var msg = this.messageToSend
          var date = new Date
          // TODO: add name selection to frontend and name to store
          var msgToSend: outgoingMessage = {
            Message: msg,
            Timestamp: date.toJSON(),
            ToIP: partner.IP,
            ToName: partner.Name
          }
          console.log(msgToSend)
          this.messageToSend = ""
          this.sendMessage(msgToSend)
        },
        check(){
          if(this.messageValid == false){
            this.messageValid = true
          }
          if(this.messageToSend == ""){
            this.messageValid = false
          }
        },
        ...mapActions([
          'sendMessage',
        ])
      },
    })

</script>