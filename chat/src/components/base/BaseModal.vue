<template>
<ion-modal ref="modal" trigger="open-modal" :presenting-element="presentingElement">
    <ion-header>
    <ion-toolbar>
        <ion-title>Neuer Chat</ion-title>
        <ion-buttons slot="end">
        <ion-button @click="dismiss()">schlie&szlig;en</ion-button>
        </ion-buttons>
    </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
    <ion-list>
      <ion-item>
        <ion-label v-if="getUserInfo().Name">Du bist: {{ getUserInfo().Name }} mit IP: {{ getUserInfo().IP }}</ion-label>
      </ion-item>
      <ion-item v-if="noPossiblePartners">
      <ion-label>
        <h2>Noch keine Partner im Netzwerk</h2>
        <p> Sobald sich jemand anmeldet, wird er hier angezeigt </p>
      </ion-label>
      </ion-item>
      <ion-item v-for="possiblePartner in getPossiblePartners()" :key="possiblePartner" @click="addNewChat(possiblePartner)">
      <ion-avatar slot="start">
          <ion-img src="https://i.pravatar.cc/300?u=b"></ion-img>
      </ion-avatar>
      <ion-label>
        <h2>{{ possiblePartner.Name }}</h2>
        <p> {{ possiblePartner.IP }} </p>
      </ion-label>
      </ion-item>
    </ion-list>
    </ion-content>
</ion-modal>
  </template>
  
  <script lang="ts">
    import {
      IonButtons,
      IonButton,
      IonModal,
      IonHeader,
      IonContent,
      IonToolbar,
      IonTitle,
      IonItem,
      IonList,
      IonAvatar,
      IonImg,
      IonLabel,
    } from '@ionic/vue';
    import { defineComponent } from 'vue';
    import { mapGetters, mapMutations } from 'vuex';
    import { User } from '../../store/index'
  
    export default defineComponent({
      components: {
        IonButtons,
        IonButton,
        IonModal,
        IonHeader,
        IonContent,
        IonToolbar,
        IonTitle,
        IonItem,
        IonList,
        IonAvatar,
        IonImg,
        IonLabel,
      },
      data() {
        return {
          presentingElement: null,
        };
      },
      computed: {
        noPossiblePartners(){
            if(this.getPossiblePartners().length == 0){
                return true
            } else {
                return false
            }
        }
      },
      methods: {
        dismiss() {
          (this.$refs.modal as any).$el.dismiss();
        },
        addNewChat(possiblePartner: User){
            this.newChat(possiblePartner)
            this.dismiss()
        },
        ...mapGetters([
            'getPossiblePartners',
            'getUserInfo'
        ]),
        ...mapMutations([
            'newChat'
        ])
      },
      mounted() {
        this.presentingElement = (this.$refs.page as any).$el;
      },
    });
  </script>