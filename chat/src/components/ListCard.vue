<template>
    <ion-label v-if="partners.length == 0">Noch keine Konversationen...</ion-label>
    <ion-card v-for="partner in partners" :key="partner.From">
      <ion-card-header>
        <ion-card-title>{{ partner.Name }}</ion-card-title>
        <ion-card-subtitle>{{ partner.IP }}</ion-card-subtitle>
      </ion-card-header>
      <ion-card-content v-if="partner.Messages">
        <ion-list>
          <ion-item>
            <form @submit.prevent="send(partner)">
            <p v-if="!messageValid">Bitte geben sie eine Nachricht ein!</p>
            <ion-label position="floating">Neue Nachricht verfassen</ion-label>
            <ion-input placeholder="Hier tippen..." v-model="messageToSend"></ion-input>
            <ion-button type="submit">Send</ion-button>
            </form>
          </ion-item>
          <ion-item-group v-for="message in partner.Messages" :key="message">

            <ion-item>
              <ion-label class="ion-text-wrap">{{ message.Message }}</ion-label>
            </ion-item>

            <ion-item-divider>
              <ion-label>
                {{ message.From + " um " + new Date(message.Timestamp).toLocaleString() }}
              </ion-label>
            </ion-item-divider>
          </ion-item-group>

        </ion-list>
        
      </ion-card-content>
    </ion-card>
  </template>
  
  <script lang="ts">
    import { IonInput, IonCard, IonCardContent, IonCardHeader, IonCardSubtitle, IonCardTitle, IonItemGroup, IonItem, IonItemDivider, IonLabel, IonList, IonButton } from '@ionic/vue';
    import { defineComponent } from 'vue';
    import { mapActions } from 'vuex';
    import { Message, User, minifiedMessage, outgoingMessage } from '../store/index'
  
    export default defineComponent({
      components: { IonInput, IonCard, IonCardContent, IonCardHeader, IonCardSubtitle, IonCardTitle, IonItemGroup, IonItem, IonItemDivider, IonLabel, IonList, IonButton},
      props: ['partners'],
      data() {
        return {
          messageToSend: "",
          messageValid: true
        }
      },
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
          'sendMessage'
        ])
      },
      created(){
        console.log(this.partners)
      }
    });
  </script>
  
  <style scoped>
    ion-item {
      --padding-start: 0;
    }
  
    /* iOS places the subtitle above the title */
    ion-card-header.ios {
      display: flex;
      flex-flow: column-reverse;
    }
  </style>