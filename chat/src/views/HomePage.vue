<template>
  <ion-page>
    <ion-header :translucent="true">
      <ion-toolbar>
        <ion-title>SchardtChat</ion-title>
        <ion-label slot="end" v-if="getUserInfo().Name">angemeldet als: {{ getUserInfo().Name }} mit IP: {{ getUserInfo().IP }}</ion-label>
      </ion-toolbar>
    </ion-header>
    <ion-content :fullscreen="true">
      <ion-header collapse="condense">
        <ion-toolbar>
        <ion-title>SchardtChat</ion-title>
        <ion-label slot="end" v-if="getUserInfo().Name">angemeldet als: {{ getUserInfo().Name }}</ion-label>
      </ion-toolbar>
    </ion-header>

    <ListCardVue :partners="getPartners()"></ListCardVue>
    
    <ion-fab slot="fixed" vertical="top" horizontal="center" :edge="true">
      <ion-fab-button>
        <ion-icon :icon="chevronDownCircle"></ion-icon>
      </ion-fab-button>
      <ion-fab-list side="bottom">
        <ion-fab-button @click="toggleDarkTheme()">
          <ion-icon :icon="contrastOutline"></ion-icon>
        </ion-fab-button>
        <ion-fab-button id="open-modal" expand="block">
            <ion-icon :icon="add"></ion-icon>
        </ion-fab-button>
      </ion-fab-list>
    </ion-fab>
  </ion-content>
  </ion-page>
</template>

<script lang="ts">
import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar, IonIcon, IonLabel, IonFab, IonFabButton, IonFabList } from '@ionic/vue';
import { star, cog, chevronDownCircle, add, contrastOutline } from 'ionicons/icons';
import { defineComponent } from 'vue';
import { mapGetters, mapActions } from 'vuex';
import ListCardVue from '@/components/ListCard.vue';
import { User } from 'src/store/index'

export default defineComponent({
  name: 'HomePage',
  components: {
    IonContent,
    IonHeader,
    IonPage,
    IonTitle,
    IonToolbar,
    IonIcon,
    ListCardVue,
    IonLabel,
    IonFab,
    IonFabButton,
    IonFabList
  },
  data(){
    return {
      darkMode: false
    }
  },
  beforeMount(){
    this.initDarkTheme()
    this.requestUserInfo()
    this.establishMessageWebsocketConnection()
    this.establishUserWebsocketConnection()
  },
  setup() {
    return { star, contrastOutline, cog, chevronDownCircle, add }
  },
  computed: {
    userInfo(): User {
      return this.getUserInfo()
    }
  },
  methods: {
    initDarkTheme(){
      const dark = localStorage.getItem('darkMode')
      if(dark == 'true'){
        document.body.classList.toggle('dark')
        this.darkMode = true
      }
    },
    toggleDarkTheme(){
      if(this.darkMode == false){
        document.body.classList.toggle('dark')
        this.darkMode = true
        localStorage.setItem('darkMode', 'true')
      } else {
        document.body.classList.toggle('dark')
        this.darkMode = false
        localStorage.setItem('darkMode', 'false')
      }
    },
  ...mapGetters([
    'getPartners',
    'getUserInfo'
  ]),
  ...mapActions(
    [
      'establishMessageWebsocketConnection',
      'establishUserWebsocketConnection',
      'requestUserInfo'
    ]
  ),
}
});
</script>

<style scoped>
#container {
  text-align: center;
  
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

#container strong {
  font-size: 20px;
  line-height: 26px;
}

#container p {
  font-size: 16px;
  line-height: 22px;
  
  color: #8c8c8c;
  
  margin: 0;
}

#container a {
  text-decoration: none;
}
</style>
