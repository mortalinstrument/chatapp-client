<template>
  <base-modal-vue></base-modal-vue>
  <ion-card v-if="noPartners">
    <ion-card-header>
      <ion-card-title>Noch keine Konversationen</ion-card-title>
      <ion-card-subtitle>Sobald ein Gespräch begonnen wurde, erscheint es hier</ion-card-subtitle>
    </ion-card-header>
  </ion-card>
  <ion-card v-for="partner in partners" :key="partner">
    <ion-card-header>
      <ion-card-title>{{ partner.Name }}</ion-card-title>
      <ion-card-subtitle>{{ partner.IP }}</ion-card-subtitle>
    </ion-card-header>
    <ion-card-content v-if="partner.Messages">
      <ion-list>
        <ion-item>
          <base-message-field-vue :partner="partner"></base-message-field-vue>
        </ion-item>
        <ion-item-group v-for="message in partner.Messages" :key="message">

          <ion-item>
            <ion-label class="ion-text-wrap">{{ message.Message }}</ion-label>
          </ion-item>

          <ion-item-divider>
            <ion-label>
              {{ message.From + " am " + new Date(message.Timestamp).toLocaleString() }}
            </ion-label>
          </ion-item-divider>
        </ion-item-group>

      </ion-list>
      
    </ion-card-content>
  </ion-card>
  </template>
  
  <script lang="ts">
    import { IonCard, IonCardContent, IonCardHeader, IonCardSubtitle, IonCardTitle, IonItemGroup, IonItem, IonItemDivider, IonLabel, IonList } from '@ionic/vue';
    import { defineComponent } from 'vue';
    import BaseModalVue from './base/BaseModal.vue';
    import BaseMessageFieldVue from './base/BaseMessageField.vue'
 
    export default defineComponent({
      components: { IonCard, IonCardContent, IonCardHeader, IonCardSubtitle, IonCardTitle, IonItemGroup, IonItem, IonItemDivider, IonLabel, IonList, BaseModalVue, BaseMessageFieldVue},
      props: ['partners'],
      computed: {
        noPartners(){
          if(this.partners.length == 0) {
            return true
          } else {
            return false
          }
        }
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