const {SlashCommandBuilder} = require('discord.js')
const axios = require('axios')
const wait = require('node:timers/promises').setTimeout;

// Fisher-Yates shuffle algorithm
function shuffle(array) {
    for (let i= array.length-1; i> 0;i--) {
        const j = Math.floor(Math.random() * (i+1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}


module.exports = {
    data: new SlashCommandBuilder()
        .setName('trivia')
        .setDescription('Get a random trivia question'),
    async execute(interaction) {

        try {
            const response = await axios.get('http://localhost:3000/api/trivias')
            const triviaData = response.data

            let combinedList = [...triviaData.incorrect_answers, triviaData.correct_answer];
            const shuffled = shuffle(combinedList)

            const options = shuffled.map((option, index) => `${index+1}. ${option}`).join('\n');            
            await interaction.reply(`Question: ${triviaData.question}\nOptions:\n${options}`)

            await wait(5_000);
            await interaction.followUp(`Answer: ${triviaData.correct_answer}`)
        } catch (error) {
           console.log(error) 
           await interaction.reply('An error occured while obtaining trivia question.')
        }
    }
}