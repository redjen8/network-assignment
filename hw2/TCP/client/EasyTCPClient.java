package hw2.TCP.client;

import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.Scanner;
import java.util.NoSuchElementException;
import java.io.*;

public class EasyTCPClient{

    private static final String SERVER_NAME = "localhost";
    private static final int SERVER_PORT = 22848;

    public static void main(String[] args) {
        
        Socket socket = null;

        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("Bye bye~"); 
        }));

        try {
            socket = new Socket();
            socket.connect(new InetSocketAddress(SERVER_NAME, SERVER_PORT));
            byte[] inputBuffer = null;
            String recvMessage = null;
            OutputStream os = socket.getOutputStream();
            InputStream is = socket.getInputStream();
            boolean loopCondition = true;
            Scanner sc = new Scanner(System.in);

            while (loopCondition) {
                System.out.println("\n<Menu>");
                System.out.println("Please input your client option in integer");
                System.out.println("Option 1) Convert text to upper case letters");
                System.out.println("Option 2) Ask the server what is my ip and port number");
                System.out.println("Option 3) Ask the server how many requests it has served so far");
                System.out.println("Option 4) Ask the server how long it has been running since it started");
                System.out.println("Option 5) Exit");
                System.out.print("Option :: ");
                
                int userOption = Integer.parseInt(sc.nextLine());
                // sc.nextLine();
                String messageCmd = null;
                int readByteCount;
                switch (userOption) {
                    case 1:
                        messageCmd = "ASK_TXTCONV";
                        System.out.print("Input lowercase sentence: ");
                        String messageText = sc.nextLine();
                        os.write((messageCmd + "," + messageText).getBytes());
                        os.flush();
                        inputBuffer = new byte[1024];
                        readByteCount = is.read(inputBuffer);
                        recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                        System.out.println("Reply from server : "+ recvMessage);
                        break;
                    case 2:
                        messageCmd = "ASK_IP_PORT";
                        os.write(messageCmd.getBytes());
                        os.flush();
                        inputBuffer = new byte[1024];
                        readByteCount = is.read(inputBuffer);
                        recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                        System.out.println("Reply from server : "+ recvMessage);
                        break;
                    case 3:
                        messageCmd = "ASK_REQ_NUM";
                        os.write(messageCmd.getBytes());
                        os.flush();
                        inputBuffer = new byte[1024];
                        readByteCount = is.read(inputBuffer);
                        recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                        System.out.println("Reply from server : "+ recvMessage);
                        break;
                    case 4:
                        messageCmd = "ASK_RUNTIME";
                        os.write(messageCmd.getBytes());
                        os.flush();
                        inputBuffer = new byte[1024];
                        readByteCount = is.read(inputBuffer);
                        recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                        System.out.println("Reply from server : requests served = "+ recvMessage);
                        break;
                    case 5:
                        messageCmd = "ASK_CONNEND";
                        os.write(messageCmd.getBytes());
                        os.flush();
                        loopCondition = false;
                        break;
                    default:
                        break;
                }
            }
            sc.close();
            is.close();
            os.close();
            socket.close();
        }
        catch (IOException e) {
            e.printStackTrace();
        }
        if (!socket.isClosed()) {
            try {
                socket.close();
            }
            catch (IOException e) {
                e.printStackTrace();
            }
        }
        // System.out.println("Bye bye~");
    }
}