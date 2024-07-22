<template>
    <div>
      <textarea ref="input" @input="handleInput"></textarea>
      <Vue3EmojiPicker v-model="selectedEmoji" @select="onEmojiSelect" />
    </div>
  </template>
  
  <script>
  import { ref } from 'vue';
  
  export default {
    setup() {
      const inputRef = ref(null);
      let selectedEmoji = ref('');
  
      function handleInput(event) {
        // 更新文本区域内容
      }
  
      function onEmojiSelect(emoji) {
        // 当用户选择emoji时，将其插入到textarea中
        if (inputRef.value) {
          inputRef.value.focus();
          const start = inputRef.value.selectionStart;
          const end = inputRef.value.selectionEnd;
          const text = inputRef.value.value;
          inputRef.value.value = `${text.slice(0, start)}${emoji}${text.slice(end)}`;
        }
        selectedEmoji.value = '';
      }
  
      return {
        inputRef,
        selectedEmoji,
        handleInput,
        onEmojiSelect
      };
    }
  };
  </script>